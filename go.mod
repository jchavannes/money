module github.com/jchavannes/money

replace github.com/jchavannes/jgo => ../jgo

go 1.16

require (
	github.com/jchavannes/jgo v0.0.0-20230307052205-37135d9964ad
	github.com/jinzhu/gorm v1.9.16
	github.com/spf13/cobra v1.1.3
	golang.org/x/crypto v0.7.0
)
