#!/usr/bin/env bash
cd ../internal/model
gen_sql_model  --ddlpath=user.sql --package=model > user_model.go && gofmt -w user_model.go
gen_sql_model  --ddlpath=sys.sql --package=model > system_model.go && gofmt -w system_model.go
