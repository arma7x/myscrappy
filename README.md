# myscrappy

> scrap data from website for personal usage

##### publicinfobanjir
```
- /publicinfobanjir/api/v1/state
- /publicinfobanjir/api/v1/river?html=0&state=KEL
- /publicinfobanjir/api/v1/rain?html=0&state=KEL
```
set html=1 for html output

##### financialtimes
```
- /ft/api/v1/currencies?group=Majors
- /ft/api/v1/commodities
- /ft/api/v1/bondsandrates
- /ft/api/v1/governmentbondsspreads
- /ft/api/v1/equities?rankingType=highestvolume&rankingSet=SP500
```
[group](https://github.com/arma7x/myscrappy/blob/master/modules/financialtimes/financialtimes.go#L3-L9)
[rangkingSet](https://github.com/arma7x/myscrappy/blob/master/modules/financialtimes/financialtimes.go#L11-L40)
[rankingType](https://github.com/arma7x/myscrappy/blob/master/modules/financialtimes/financialtimes.go#L42-L46)
