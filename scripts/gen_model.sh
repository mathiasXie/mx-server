#!/usr/bin/env bash
cd ../applications/xiaozhi-server/internal/model
gen_sql_model  --ddlpath=ai_devices.sql --package=model > ai_devices_model.go && gofmt -w ai_devices_model.go
gen_sql_model  --ddlpath=ai_users.sql --package=model > ai_users_model.go && gofmt -w ai_users_model.go
gen_sql_model  --ddlpath=ai_messages.sql --package=model > ai_messages_model.go && gofmt -w ai_messages_model.go
gen_sql_model  --ddlpath=ai_roles.sql --package=model > ai_roles_model.go && gofmt -w ai_roles_model.go
