## go-aws-migrator

This repository contains a migrator package written in go which helps run and manage database migrations. State of
migrations is saving in
the [AWS DynamoDB](https://aws.amazon.com/dynamodb/?trk=ps_a134p000006padwAAA&trkCampaign=acq_paid_search_brand&sc_channel=PS&sc_campaign=acquisition_EEM&sc_publisher=Google&sc_category=Database&sc_country=EEM&sc_geo=EMEA&sc_outcome=acq&sc_detail=dynamodb&sc_content=DynamoDB_e&sc_matchtype=e&sc_segment=536393757514&sc_medium=ACQ-P|PS-GO|Brand|Desktop|SU|Database|DynamoDB|EEM|EN|Text|xx|EU&s_kwcid=AL!4422!3!536393757514!e!!g!!dynamodb&ef_id=CjwKCAiA6seQBhAfEiwAvPqu11oruz5QdnON14qU1_71ZxjKCxGS1saSrZqab-ourmxi7NH4a1-J_BoCmasQAvD_BwE:G:s&s_kwcid=AL!4422!3!536393757514!e!!g!!dynamodb)
. In `src/examples` you have working example for entire concept. In shortcut, you have to write a set of migration
definitions:

```go
type Definition struct {
Name string
Func func () error
}
```

Pass them to `Run` function as argument. Next, this function iterates through your all definitions set and execute your
code for every, single definition. Information about all previous migration are saved in a DynamoDB database. When you
try to run the migration again, only new definitions will be fired.