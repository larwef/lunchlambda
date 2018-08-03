# lunchLambda
Program designed to run on AWS Lambda. Gets Menu from a **Getter** and transforms it to a **Menu** object.
The Menu object can be modified by a **Modifier** before being sent by a **Sender**.

The **Runner** object is used to drive the application. A **Runner** can be configured with one **Getter** to
get a **Menu** object, n modifier and n senders. The **Menu** object will be modified in the order the modifiers
are added, and sent in the same order the **Senders** are added. If more than one **Getter** are needed just
create more **Runner**s. 

## Getters
Objects implementing the Getter interface defined by **Menu**. Used to get Menues from a source and return
a Menu object.
### Implemented Getters
#### Braathens
Webscraper getting Menu iterms from http://www.braatheneiendom.no/NO/Eiendommer/DEG16/kjokkenet/.

## Modifiers
Objects implementing the Modifier interface defined by **Menu**. Used to modify Menu objects before **Sender**
is invoked.
### Implemented Modifiers
#### Audio
Uses Amazon Polly to create an .mp3 file of the Menu item and puts the file in an S3 bucket.

## Senders
Objects implementing the Sender interface defined by **Menu**. Used to send a formatted message to an endpoint.
### Implemented Senders
#### Slack
Posts the Menu to a Slack webhook.
#### SNS (AWS)
Puts the Menu on a SNS topic.