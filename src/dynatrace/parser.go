package dynatrace

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

/*
Sample Visit
{"visitorId":"954","visitId":"954+774","tenantId":"bsu97947","startTime":1504071473250,"endTime":1504071544606,"visitType":"REAL_USER","ip":"99.109.183.90","geoInfos":[{"name":"North America","kind":"CONTINENT"},{"name":"United States","kind":"COUNTRY"}],"osInfos":[{"name":"Android","kind":"FAMILY"},{"name":"Android 6.0.1","kind":"VERSION"}],"browserInfos":[],"deviceInfos":[],"userActions":[{"name":"AppStart (easyTravel)","startTime":1504071473250,"endTime":1504071473273,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Display SearchJourneyActivity","startTime":1504071492188,"endTime":1504071492214,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on User Account","startTime":1504071502897,"endTime":1504071502906,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on Search","startTime":1504071513734,"endTime":1504071513743,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on Paris","startTime":1504071517572,"endTime":1504071517581,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on New York","startTime":1504071521581,"endTime":1504071521588,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on Search","startTime":1504071525441,"endTime":1504071525449,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on Berlin","startTime":1504071529324,"endTime":1504071529332,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on Amsterdam","startTime":1504071533231,"endTime":1504071533236,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on Vienna","startTime":1504071537006,"endTime":1504071537014,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on Search","startTime":1504071540753,"endTime":1504071540759,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0},{"name":"Touch on Rome","startTime":1504071544601,"endTime":1504071544606,"type":"Custom","application":"easyTravel Demo","apdexType":"UNKNOWN","networkTime":0,"serverTime":0,"frontendTime":0,"documentInteractiveTime":0,"jsErrors":0,"failedImages":0,"failedXhrRequests":0,"httpRequestsWithErrors":0,"thirdPartyResources":0,"thirdPartyBusyTime":0,"cdnResources":0,"cdnBusyTime":0,"firstPartyResources":0,"firstPartyBusyTime":0}],"mobile":true,"isp":"ATT-INTERNET4 - AT&T Services, Inc.","sessionTag":"bayu","clientType":"Mobile App","stringTags":{},"numTags":{},"dateTags":{},"newVisitor":true,"screenHeight":1080,"screenWidth":1920,"clientHeight":1080,"clientWidth":1920,"manufacturer":"Sony","deviceName":"Xperia Z5","internalSynthetic":false}
*/

type KeyValue struct {
	Name				string				`json:"name"`
	Kind				string				`json:"kind"`
}

type UserAction struct {
	Name				string				`json:"name"`
	StartTime			uint64				`json:"startTime"`
	EndTime				uint64				`json:"endTime"`
	Type				string				`json:"type"`
	Application			string				`json:"application"`
}

func (action UserAction) String() string {
	return action.Name
}

type Visit struct {
	VisitorID			string				`json:"visitorId"`
	VisitID				string				`json:"visitId"`
	TenantID			string				`json:"tenantId"`
	StartTime			uint64				`json:"startTime"`
	EndTime				uint64				`json:"endTime"`
	VisitType			string				`json:"visitType"`
	IP					string				`json:"ip"`
	GeoInfos			[]KeyValue			`json:"geoInfos"`
	OSInfos				[]KeyValue			`json:"osInfos"`
	BrowserInfos		[]KeyValue			`json:"browserInfos"`
	DeviceInfos			[]KeyValue			`json:"deviceInfos"`
	UserActions			[]UserAction		`json:"userActions"`
	Mobile				bool				`json:"mobile"`
	ISP					string				`json:"isp"`
	SessionTag			string				`json:"sessionTag"`
	ClientType			string				`json:"clientType"`
	NewVisitor			bool				`json:"newVisitor"`
	ScreenHeight		int					`json:"screenHeight"`
	ScreenWidth			int					`json:"screenWidth"`
	ClientHeight		int					`json:"clientHeight"`
	ClientWidth			int					`json:"clientWidth"`
	InternalSynthetic	bool				`json:"internalSynthetic"`
}

func (v Visit) String() string {
	var b bytes.Buffer

	b.WriteString("VisitID: ")
	b.WriteString(v.VisitID)
	b.WriteString("\n")
	b.WriteString("Session Tag: ")
	b.WriteString(v.SessionTag)
	b.WriteString("\n")
	b.WriteString("User Actions:\n")
	for _, action := range v.UserActions {
		b.WriteString("\t")
		b.WriteString(action.String())
		b.WriteString("\n")
	}

	return b.String()
}


func Parse(r io.Reader) []Visit {
	var visits []Visit

	dec	:= json.NewDecoder(r)
	for dec.More() {
		var visit Visit
		err := dec.Decode(&visit)
		if err != nil {
			log.Println("Error parsing visit:", err)
			continue
		}
		visits = append(visits, visit)
	}

	return visits
}
