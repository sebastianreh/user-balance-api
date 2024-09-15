Possible improvements:

Add job handling for uploading transactions but since the challenge states: "For every inconsistent input data such as user not found or bad datetime format, 
a 400 bad request response must be returned." this is not possible. A solution would be to add a job entity that supports the job execution in background
and then look for the job status to see if there is some error.

Also in the challenge there is an inconsistency in this request:

"/users/{user_id}/balance?from=YYYY-MM-DDThh:mm:ssZ&to=YYYY-MM-DDThh:mm:ss"

Here, the from date is with timezone and the to date is without timezone. This implementation doesn't support timezones, since it's the common
practice to use GMT in this type of services.

I added users and transactions endpoint so if there is any need for modification connecting to the database will not be necessary. 
- Add unique user validation for creation