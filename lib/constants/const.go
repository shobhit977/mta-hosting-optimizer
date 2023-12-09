package constants

var (
	Bucket       = "mta-hosting-bucket" //bucket name must be unique. Change this value if you deploy your code
	Key          = "ipConfig.json"
	ThresholdKey = "threshold" // environment variable is stored in lambda
	Region       = "ap-south-1"
)
