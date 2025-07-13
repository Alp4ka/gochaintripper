module github.com/Alp4ka/gochaintripper/logging/json/tidwall

go 1.23.0

replace (
	github.com/Alp4ka/gochaintripper => ./../../../
	github.com/Alp4ka/gochaintripper/logging => ./../../
)

require (
	github.com/Alp4ka/gochaintripper/logging v0.0.0-local
	github.com/tidwall/gjson v1.18.0
	github.com/tidwall/sjson v1.2.5
)

require (
	github.com/Alp4ka/gochaintripper v0.0.0-local // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
)
