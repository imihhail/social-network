package database

import (
	"backend/database/sqlite"
	"backend/helpers"
	"backend/structs"
	"database/sql"
	"fmt"
)

// is email in table
func GetEmailIfExists(email string) bool {
	db := sqlite.DbConnection()
	var userId string
	command := "SELECT id FROM users WHERE email=?"
	err := db.QueryRow(command, email).Scan(&userId)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("GetEmailIfExists", err)
		}
		return false
	}
	defer db.Close()
	return true
}

// is username in table
func GetUsernameIfExists(username string) bool {
	db := sqlite.DbConnection()
	var userId string
	command := "SELECT id FROM users WHERE username=?"
	err := db.QueryRow(command, username).Scan(&userId)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("GetUsernameIfExists", err)
		}
		return false
	}
	defer db.Close()
	return true
}

// get userid and password by email
func GetUserIdPswByEmail(email string) (string, string) {
	db := sqlite.DbConnection()
	var userId string
	var userPsw string
	command := "SELECT id, password FROM users WHERE email=?"
	err := db.QueryRow(command, email).Scan(&userId, &userPsw)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("GetUserIdPswByEmail", err)
		}
		return "0", "Email error"
	}
	defer db.Close()
	return userId, userPsw
}

// get userid by email
func GetUserIdByEmail(email string) string {
	db := sqlite.DbConnection()
	var userId string
	command := "SELECT id FROM users WHERE email=?"
	err := db.QueryRow(command, email).Scan(&userId)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("GetUserIdByEmail", err)
		}
		return "0"
	}
	defer db.Close()
	return userId
}

// get userId session by hash from session table
func GetUserSession(cookie string) string {
	db := sqlite.DbConnection()
	command := "SELECT user_fk_users FROM session WHERE hash=?"
	err := db.QueryRow(command, cookie).Scan(&cookie)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("GetUserSession", err)
		}
		return "0"
	}
	defer db.Close()
	return cookie
}

func GetUserProfile(userId string) structs.Profile {
	db := sqlite.DbConnection()
	var userProfile structs.Profile
	command := "SELECT id, email, first_name, last_name, date_of_birth, username, about_user, avatar FROM users WHERE id=?"
	err := db.QueryRow(command, userId).Scan(&userProfile.Id,
		&userProfile.Email,
		&userProfile.FirstName,
		&userProfile.LastName,
		&userProfile.DateOfBirth,
		&userProfile.Username,
		&userProfile.AboutUser,
		&userProfile.Avatar)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("getUserProfile", err)
		}
		fmt.Println("User email not found in users table!")
	}
	defer db.Close()
	return userProfile
}

func GetProfilePosts(userId string) []structs.ProfilePosts {
	db := sqlite.DbConnection()
	var profilePosts []structs.ProfilePosts
	command := "SELECT id, user_fk_users, post_content, post_image, privacy_fk_posts_privacy, date FROM posts WHERE user_fk_users=?"
	rows, err := db.Query(command, userId)
	if err != nil {
		helpers.CheckErr("GetProfilePosts", err)
	}
	for rows.Next() {
		var profilePost structs.ProfilePosts
		err = rows.Scan(&profilePost.PostId,
			&profilePost.UserId,
			&profilePost.PostContent,
			&profilePost.Picture,
			&profilePost.PostPrivacy,
			&profilePost.Date)
		if err != nil {
			helpers.CheckErr("GetProfilePosts", err)
		}
		profilePosts = append(profilePosts, profilePost)
	}
	defer db.Close()
	return profilePosts
}

func GetAllPosts() []structs.Posts {
	db := sqlite.DbConnection()
	defer db.Close()

	var allPosts []structs.Posts

	command := "SELECT posts.id, users.username, users.avatar, posts.post_Title, posts.post_content, posts.post_image, posts.privacy_fk_posts_privacy, posts.date, users.email FROM posts INNER JOIN users ON posts.user_fk_users == users.id ORDER BY posts.date DESC"
	rows, err := db.Query(command)
	if err != nil {
		helpers.CheckErr("getAllPosts", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var post structs.Posts
		err = rows.Scan(&post.PostID, &post.Username, &post.Avatar, &post.Title, &post.Content, &post.Picture, &post.Privacy, &post.Date, &post.Email)
		if err != nil {
			helpers.CheckErr("getAllPosts", err)
			continue
		}
		allPosts = append(allPosts, post)
	}

	if err = rows.Err(); err != nil {
		helpers.CheckErr("getAllPosts", err)
	}
	return allPosts
}

func GetNewPost() structs.Posts {
	db := sqlite.DbConnection()
	var lastPost structs.Posts

	command := "SELECT posts.id, users.username, users.avatar, posts.post_title, posts.post_content, posts.post_image, posts.privacy_fk_posts_privacy, posts.date FROM posts INNER JOIN users ON posts.user_fk_users == users.id ORDER BY posts.date DESC LIMIT 1"
	err := db.QueryRow(command).Scan(&lastPost.PostID, &lastPost.Username, &lastPost.Avatar, &lastPost.Title, &lastPost.Content, &lastPost.Picture, &lastPost.Privacy, &lastPost.Date)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("getNewpost", err)
		}
		fmt.Println("Error selecting last post")
	}
	defer db.Close()
	return lastPost
}

func GetAllUsers(userId string) []structs.Profile {
	db := sqlite.DbConnection()
	var allUsers []structs.Profile

	command := "SELECT id, username, email FROM users"
	rows, err := db.Query(command)
	if err != nil {
		helpers.CheckErr("getAllPosts", err)
		return nil
	}

	for rows.Next() {
		var user structs.Profile
		err = rows.Scan(&user.Id, &user.Username, &user.Email)
		if err != nil {
			helpers.CheckErr("getAllPosts", err)
			continue
		}
		if user.Id != userId {
			allUsers = append(allUsers, user)
		}
	}

	defer rows.Close()

	defer db.Close()
	return allUsers
}

func GetAllComments(postID string) []structs.Comments {
	db := sqlite.DbConnection()
	defer db.Close()

	var allComments []structs.Comments

	command := "SELECT users.username, users.avatar, comments.comment_content, comments.comment_image, comments.date FROM comments INNER JOIN users ON user_fk_users == users.id where post_Id_fk_posts = ? ORDER BY comments.date DESC"
	rows, err := db.Query(command, postID)
	if err != nil {
		helpers.CheckErr("getAllComments", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var comment structs.Comments
		err = rows.Scan(&comment.Username, &comment.Avatar, &comment.Content, &comment.Picture, &comment.Date)
		if err != nil {
			helpers.CheckErr("getAllComments loop: ", err)
			continue
		}
		allComments = append(allComments, comment)
	}

	if err = rows.Err(); err != nil {
		helpers.CheckErr("getAllComments rows: ", err)
	}
	return allComments
}

func GetNewComment() structs.Comments {
	db := sqlite.DbConnection()
	var newComment structs.Comments

	command := "SELECT users.username, users.avatar, comments.comment_content, comments.comment_image, comments.date FROM comments INNER JOIN users ON user_fk_users == users.id ORDER BY comments.date DESC LIMIT 1"
	err := db.QueryRow(command).Scan(&newComment.Username, &newComment.Avatar, &newComment.Content, &newComment.Picture, &newComment.Date)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("getNewComment", err)
		}
		fmt.Println("Error selecting new comment")
	}
	defer db.Close()
	return newComment
}

// Get user privacy setting
func GetUserPrivacy(userId string) string {
	db := sqlite.DbConnection()
	var privacy string
	command := "SELECT privacy_fk_users_privacy FROM user_privacy WHERE user_fk_users=?"
	err := db.QueryRow(command, userId).Scan(&privacy)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("GetUserPrivacy", err)
		}
		return "0"
	}
	defer db.Close()
	return privacy
	// 1 = public, 2 = private, (3 = almost private)???
}

// get userid if email in table
func GetUserIdIfEmailExists(email string) string {
	db := sqlite.DbConnection()
	var userId string
	command := "SELECT id FROM users WHERE email=?"
	err := db.QueryRow(command, email).Scan(&userId)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("GetEmailIfExists", err)
		}
		return "no such email"
	}
	defer db.Close()
	return userId
}

func GetMessages(fromuser, touser string) []structs.ChatMessage {
	db := sqlite.DbConnection()
	var UserMessages []structs.ChatMessage
	command := "SELECT * FROM messages WHERE (message_sender_fk_users = ? AND message_receiver_fk_users = ?) OR (message_receiver_fk_users = ? AND message_sender_fk_users = ?) ORDER BY id"
	rows, err := db.Query(command, touser, fromuser, touser, fromuser)
	if err != sql.ErrNoRows {
		helpers.CheckErr("GetMessages", err)
	}
	defer db.Close()
	helpers.CheckErr("GetMessage", err)
	for rows.Next() {
		var message structs.ChatMessage
		rows.Scan(&message.ChatMessageId, &message.MessageSender, &message.Message, &message.MessageReceiver, &message.Date)
		UserMessages = append(UserMessages, message)
	}
	defer rows.Close()
	return UserMessages
}

func GetAllGroups() []structs.Groups {
	db := sqlite.DbConnection()
	defer db.Close()

	var allGroups []structs.Groups

	command := `SELECT guilds.id, guilds.creator_fk_users, users.email, guilds.guild_title, guilds.guild_description, guilds.date
				FROM guilds 
				INNER JOIN users ON guilds.creator_fk_users == users.id 
				GROUP BY guilds.id
				ORDER BY guilds.date DESC`

	rows, err := db.Query(command)
	if err != nil {
		helpers.CheckErr("Selecting getAllGroups", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var group structs.Groups
		err = rows.Scan(&group.Id, &group.CreatorId, &group.Creator, &group.Title, &group.Description, &group.Date)
		if err != nil {
			helpers.CheckErr("Iterating GetAllGroups", err)
			continue
		}
		allGroups = append(allGroups, group)
	}

	if err = rows.Err(); err != nil {
		helpers.CheckErr("GetAllGroups rows", err)
	}
	return allGroups
}

func GetNewGroup() structs.NewGroup {
	db := sqlite.DbConnection()
	var newGroup structs.NewGroup

	command := "SELECT guilds.id, guilds.creator_fk_users, users.username, users.id, guilds.guild_title, guilds.guild_description, guilds.date FROM guilds INNER JOIN users ON guilds.creator_fk_users == users.id ORDER BY guilds.date DESC LIMIT 1"
	err := db.QueryRow(command).Scan(&newGroup.Id, &newGroup.CreatorId, &newGroup.Creator, &newGroup.Members, &newGroup.Title, &newGroup.Description, &newGroup.Date)
	if err != nil {
		if err != sql.ErrNoRows {
			helpers.CheckErr("getNewGroup", err)
		}
		fmt.Println("Error selecting new group")
	}
	defer db.Close()
	return newGroup
}

func GetNotifications(currentUser string) []structs.GrInvNotifications {
	db := sqlite.DbConnection()
	defer db.Close()

	var allNotif []structs.GrInvNotifications

	command := `SELECT sender_fk_users, reciever_fk_users, email, guild_title, guildnotifications.id, guilds.id FROM guildnotifications
				INNER JOIN guilds ON guildnotifications.guildid_fk_guilds = guilds.id
				INNER JOIN users ON guildnotifications.sender_fk_users = users.id
				WHERE reciever_fk_users = ?`

	rows, err := db.Query(command, currentUser)
	if err != nil {
		helpers.CheckErr("GetNotifications selecting error: ", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var notif structs.GrInvNotifications

		err = rows.Scan(&notif.SenderId, &notif.RecieverId, &notif.SenderEmail, &notif.GroupTitle, &notif.NotificationId, &notif.GroupId)
		if err != nil {
			helpers.CheckErr("GetNotifications Next error: ", err)
			continue
		}
		allNotif = append(allNotif, notif)
	}

	if err = rows.Err(); err != nil {
		helpers.CheckErr("GetNotifications", err)
	}
	return allNotif
}

func GetEventNotifications(currentUser string) []structs.EventNotifications {
	db := sqlite.DbConnection()
	defer db.Close()

	var allNotif []structs.EventNotifications

	command := `SELECT event_notifications.id, sender_fk_users, reciever_fk_users, email, event_title, event_description, event_time, events.id, guild_id_fk_guilds, guild_title 
				FROM event_notifications
				INNER JOIN events ON event_notifications.event_id_fk_events = events.id
				INNER JOIN users ON event_notifications.sender_fk_users = users.id
				INNER JOIN guilds ON event_notifications.guild_id_fk_guilds = guilds.id
				WHERE reciever_fk_users = ?`

	rows, err := db.Query(command, currentUser)
	if err != nil {
		helpers.CheckErr("GetEventNotifications selecting error: ", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var notif structs.EventNotifications

		err = rows.Scan(&notif.NotificationId, &notif.SenderId, &notif.RecieverId, &notif.SenderEmail, &notif.EventTitle, &notif.EventDescription, &notif.EventTime, &notif.EventId, &notif.GroupId, &notif.GroupTitle)
		if err != nil {
			helpers.CheckErr("GetEventNotifications Next error: ", err)
			continue
		}
		allNotif = append(allNotif, notif)
	}

	if err = rows.Err(); err != nil {
		helpers.CheckErr("GetEventNotifications", err)
	}
	return allNotif
}

func GetGroupMembers(groupId string) []structs.GroupMember {
	db := sqlite.DbConnection()
	defer db.Close()

	var groupMembers []structs.GroupMember

	command := "SELECT members_fk_users, users.email FROM guildmembers INNER JOIN users ON members_fk_users  = users.id  WHERE guild_id_fk_guilds = ?"
	rows, err := db.Query(command, groupId)
	if err != nil {
		helpers.CheckErr("GetGroupMembers selecting error: ", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var member structs.GroupMember

		err = rows.Scan(&member.Id, &member.Email)
		if err != nil {
			helpers.CheckErr("GetGroupMembers Next error: ", err)
			continue
		}
		groupMembers = append(groupMembers, member)
	}

	if err = rows.Err(); err != nil {
		helpers.CheckErr("GetGroupMembers", err)
	}
	return groupMembers
}

func GetGroupEvents(groupId string) []structs.Events {
	db := sqlite.DbConnection()
	defer db.Close()

	var groupEvents []structs.Events

	command := `SELECT events.*, users.email 
				FROM events 
				INNER JOIN users ON events.creator_fk_users = users.id  
				WHERE guildid_fk_guilds = ?`

	rows, err := db.Query(command, groupId)
	if err != nil {
		helpers.CheckErr("GetGroupEvents select: ", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var event structs.Events

		err = rows.Scan(&event.EventId, &event.EventCreatorId, &event.GroupId, &event.EventTitle, &event.EventDescription, &event.EventTime, &event.CreatorEmail)
		if err != nil {
			helpers.CheckErr("GetGroupEvents loop: ", err)
			continue
		}

		command := `SELECT participant_fk_users, users.email FROM event_participants INNER JOIN users ON event_participants.participant_fk_users = users.id WHERE event_participants.event_id_fk_events = ?`
		rows, err := db.Query(command, event.EventId)
		if err != nil {
			helpers.CheckErr("GetGroupEvents selecting participants: ", err)
			return nil
		}
		defer rows.Close()

		for rows.Next() {
			var participant structs.EventParticipant

			err = rows.Scan(&participant.ParticipantId, &participant.ParticipantEmail)
			if err != nil {
				helpers.CheckErr("GetGroupEvents participants loop: ", err)
				continue
			}
			event.Participants = append(event.Participants, participant)
		}
		groupEvents = append(groupEvents, event)
	}

	if err = rows.Err(); err != nil {
		helpers.CheckErr("GetGroupEvents rows: ", err)
	}
	return groupEvents
}

// this is for finding out if a logged in user
// is looking at their own profile or someone else's
// I compare session owner email to /profile/ path
func GetEmailFromSession(session string) string {
	fmt.Println("Session original: ", session)
	db := sqlite.DbConnection()
	defer db.Close()

	var email string
	err := db.QueryRow("SELECT users.email FROM session INNER JOIN users ON session.user_fk_users = users.id WHERE session.hash = ?", session).Scan(&email)
	if err != nil {
		helpers.CheckErr("GetEmailFromSession", err)
		return ""
	}

	return email
}
