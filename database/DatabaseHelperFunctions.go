package database

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	error2 "kontest-user-service/error"
	"kontest-user-service/model"
)

func FindUserByID(uid uuid.UUID) (*model.User, error) {
	var user model.User

	query := `SELECT * FROM user_info WHERE id = :id`

	// Use NamedQuery to retrieve the user by ID
	rows, err := GetDB().NamedQuery(query, map[string]interface{}{
		"id": uid,
	})

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Check if any rows are returned
	if rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, fmt.Errorf("error scanning result: %v", err)
		}

		userSites, _ := getSitesOfAUser(uid)

		if userSites == nil {
			userSites = []model.Site{}
		}

		user.Sites = userSites

		return &user, nil
	}

	// If no rows were returned, the user was not found
	return nil, &error2.UserNotFoundError{}
}

func getSitesOfAUser(uid uuid.UUID) ([]model.Site, error) {
	// Second query: Get associated sites for the user
	querySites := `SELECT site_name, is_site_enabled, is_automatic_calendar_notification_enabled, seconds_before_which_app_notification_to_set FROM user_site_info WHERE user_id = $1`

	var sites []model.Site

	rows, err := GetDB().Query(querySites, uid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var site model.Site

		// Scan the fields, using pq.Array for the array type
		var secondsBefore pq.Int32Array // Temporary variable for the array
		if err := rows.Scan(&site.SiteName, &site.IsSiteEnabled, &site.IsAutomaticCalendarNotificationEnabled,
			&secondsBefore); err != nil {
			return nil, err
		}

		// Convert pq.Int32Array to []int
		site.SecondsBeforeWhichAppNotificationToSet = make([]int, len(secondsBefore))
		for i, v := range secondsBefore {
			site.SecondsBeforeWhichAppNotificationToSet[i] = int(v)
		}

		// Append the site to the slice
		sites = append(sites, site)
	}

	return sites, nil
}

func UpdateUserOrCreate(user *model.User, tx *sqlx.Tx) (bool, error) {
	query := `
	INSERT INTO user_info (
		id, first_name, last_name, college_name, college_state,
	    leetcode_username, codechef_username, codeforces_username, 
		min_duration_in_seconds, max_duration_in_seconds, account_create_date
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT (id) DO UPDATE
	SET
		first_name = EXCLUDED.first_name,
		last_name = EXCLUDED.last_name,
		college_name = EXCLUDED.college_name,
		college_state = EXCLUDED.college_state,
		leetcode_username = EXCLUDED.leetcode_username,
		codechef_username = EXCLUDED.codechef_username,
		codeforces_username = EXCLUDED.codeforces_username,
		min_duration_in_seconds = EXCLUDED.min_duration_in_seconds,
		max_duration_in_seconds = EXCLUDED.max_duration_in_seconds;
	`

	var err error

	if tx == nil {
		tx = db.MustBegin()

		// Defer rollback unless the transaction is successfully committed
		defer func() {
			if r := recover(); r != nil {
				_ = tx.Rollback()
				panic(r) // Re-panic after rollback if a panic occurred
			} else if err != nil {
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}()
	}

	// Execute the upsert query with user data
	_, err = tx.Exec(query,
		user.ID, // Assuming `user.ID` is the primary key in user_info
		user.FirstName,
		user.LastName,
		user.CollegeName,
		user.CollegeState,
		user.LeetcodeUsername,
		user.CodechefUsername,
		user.CodeforcesUsername,
		user.MinDurationInSecond,
		user.MaxDurationInSecond,
		user.AccountCreateDate,
	)

	if user.Sites == nil {
		user.Sites = []model.Site{}
	}

	_, err = updateSites(user.ID, user.Sites, tx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func updateSites(uid uuid.UUID, sites []model.Site, tx *sqlx.Tx) (bool, error) {
	query := `
INSERT INTO user_site_info (site_name, is_site_enabled, is_automatic_calendar_notification_enabled, seconds_before_which_app_notification_to_set, user_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, site_name) DO UPDATE SET
	site_name = EXCLUDED.site_name,
	is_site_enabled = EXCLUDED.is_site_enabled,
	is_automatic_calendar_notification_enabled = EXCLUDED.is_automatic_calendar_notification_enabled,
	seconds_before_which_app_notification_to_set = EXCLUDED.seconds_before_which_app_notification_to_set
`

	var err error

	if tx == nil {
		tx = GetDB().MustBegin()

		// Defer rollback unless the transaction is successfully committed
		defer func() {
			if r := recover(); r != nil {
				_ = tx.Rollback()
				panic(r) // Re-panic after rollback if a panic occurred
			} else if err != nil {
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}()
	}

	stmt, err := tx.Preparex(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	for _, site := range sites {
		_, err := stmt.Exec(
			site.SiteName,
			site.IsSiteEnabled,
			site.IsAutomaticCalendarNotificationEnabled,
			pq.Array(site.SecondsBeforeWhichAppNotificationToSet),
			uid,
		)
		if err != nil {
			return false, err
		}

	}

	return true, nil
}
