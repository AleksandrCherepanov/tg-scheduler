Table notification.telegram_user {
  id int [pk]
  name varchar
  created_at timestamp
}

Table notification.notification {
  id int [pk, increment]
  user_id int [ref: > notification.telegram_user.id]
  value text
  schedule varchar
  created_at timestamp
  updated_at timestamp
  
  Indexes {
    (user_id) [name: 'user_id']
  }
}

