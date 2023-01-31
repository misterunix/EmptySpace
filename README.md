# emptyspace
Realtime 4x Space Game

## Dont expect fast progress.

1. The database / reflect code works. Time to start thinking about mechanics.
2. PAUSE to ask for advice on reddit. 

## A.I.

A.I. runs on a fixed clock cycle. It would be to easy for the A.I. to outplay the human.

The next calculated event is encoded and inserted into the DB.  
If or when that time is reached the effect(s) are generated and actions taken. 
If this causes another event, calculate an offset and store the time.
If a time in the DB is in the event, do the actions and do a partial result.

Example, 8 ore an hour but the DB even is past. Calculate how much has been mined and then set a new event time.

Planets can only be doing one thing at a time; exept launch fleets and such. 
Automated items happen on their own.

This would be a good place to try out the time code.