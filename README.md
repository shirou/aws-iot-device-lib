# AWS IoT Device Lib

This Go library is for use on the device side of AWS IoT, built on top of the Paho MQTT library.

Note: This library is *NOT* official AWS.

Currently implemented API.

- AWS IoT Jobs

Go 1.18 or later version is required because of generics.

## AWS IoT Jobs

The API is implemented according to [AWS IoT Jobs device MQTT API](https://docs.aws.amazon.com/iot/latest/developerguide/jobs-mqtt-api.html).

### example: DescribeJobExecution

Here is an example. For examples of other APIs, see the directory under examples/jobs.

```go
// Here `mc` is a Paho mqtt.client that has already been set up.

client, err := jobs.NewClient(mc) // setup jobs.Client based on mqtt.Client
if err != nil {
	return err
}

// Create a request
req := DescribeJobExecutionInput{
	ThingName: aws.String("thing-1234"),
	JobId:     aws.String("test-job"),
}

// Calls a method as synchronous execution.
ret, err := client.DescribeJobExecution(context.Background(), req)
if err != nil {
    // If rejected, error will be returned.
	return err
}
// Now you get
for _, step := range ret.Execution.JobDocument.Steps {
	fmt.Printf("Steps: %s\n", step.Action.Name)
}
```


## License

Apache License 2.0

Some of codes are BSD-3-Clause beceause those are copied from Go code.
