module github.com/omcrgnt/builder

go 1.26.2

retract (
	[v1.0.0, v1.20.0]
	[v0.1.0, v0.20.0]
	v0.20.1 // incompatible with res v0.20.2
)

require github.com/omcrgnt/res v0.20.2
