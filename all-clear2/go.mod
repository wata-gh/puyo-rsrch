module github.com/wata-gh/puyo-rsrch/all-clear2

go 1.18

require (
	github.com/wata-gh/puyo-rsrch/common v0.0.0-20220606132353-8536d11daaf9
	github.com/wata-gh/puyo2 v0.0.0-20220606061942-c88c9fef8840
)

require github.com/BurntSushi/toml v1.1.0 // indirect

replace github.com/wata-gh/puyo2 => ../../puyo2

replace github.com/wata-gh/puyo-rsrch/common => ../common
