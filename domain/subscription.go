package domain

import "time"

type Subscription struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	// ID of subreddit
	Subreddit    string    `db:"subreddit"`
	SubscribedAt time.Time `db:"subscribed_at"`
}

/* Schema
CREATE TABLE IF NOT EXISTS subscriptions (
	id INT AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(255) NOT NULL,
	subreddit VARCHAR(255) NOT NULL,
	subscribed_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER TABLE subscriptions ADD INDEX subscriptions_username_index (username);
 */
