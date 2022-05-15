module github.com/wata-gh/puyo-rsrch/gtr-back2

go 1.18

replace github.com/wata-gh/puyo-rsrch/common => ../common

require (
	github.com/wata-gh/puyo-rsrch/common v0.0.0-00010101000000-000000000000
	github.com/wata-gh/puyo2 v0.0.0-20220511163044-07904a35cd9b
	gonum.org/v1/gonum v0.11.0
)

require github.com/BurntSushi/toml v1.1.0 // indirect

replace github.com/wata-gh/puyo2 => ../../puyo2
