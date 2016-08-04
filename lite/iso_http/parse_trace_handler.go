package iso_http

import (
	"net/http"
	"bufio"
	"log"
	"bytes"
	"strconv"
	"github.com/rkbalgi/isosim/lite/spec"
	"io/ioutil"
	"encoding/hex"
	"encoding/json"
	"github.com/rkbalgi/isosim/lite/ui_data"
)

func ParseTraceHandler() {

	http.HandleFunc(ParseTraceUrl, func(rw http.ResponseWriter, req *http.Request) {

		reqUri := req.RequestURI
		scanner := bufio.NewScanner(bytes.NewBufferString(reqUri))
		scanner.Split(splitByFwdSlash)
		urlComponents := make([]string, 0, 10);
		for scanner.Scan() {
			if (len(scanner.Text()) != 0) {
				urlComponents = append(urlComponents, scanner.Text())
			}
		}

		log.Print(urlComponents)

		if (len(urlComponents) != 5) {
			sendError(rw, "invalid url - " + reqUri);
			return;
		}

		rw.WriteHeader(200)
		paramSpecId := urlComponents[3];
		paramMsgId := urlComponents[4];

		specId, err := strconv.ParseInt(paramSpecId, 10, 0);
		if (err != nil) {
			sendError(rw, "invalid spec id in url - " + reqUri);
			return;
		}
		msgId, err := strconv.ParseInt(paramMsgId, 10, 0);
		if (err != nil) {
			sendError(rw, "invalid msg id in url - " + reqUri);
			return;
		}

		isoSpec := spec.GetSpec(int(specId));
		if (isoSpec != nil) {
			msg := isoSpec.GetMessageById(int(msgId));
			if (msg != nil) {
				log.Printf("Fetching Template for Spec: %s and Message: %s", isoSpec.Name, msg.Name);
				//TODO::
				reqData, err := ioutil.ReadAll(req.Body);
				if (err != nil) {
					sendError(rw, err.Error());
					return;
				}
				if (spec.DebugEnabled) {
					log.Print("Trace = " + string(reqData));
					msgData, err := hex.DecodeString(string(reqData));
					if (err != nil) {
						sendError(rw, " invalid trace - should be hex string " + err.Error());
						return;
					} else {
						parsedMsg, err := msg.Parse(msgData);
						if (err != nil) {

							sendError(rw, "parse error " + err.Error());
							return;
						}

						fieldDataList := ToJsonList(parsedMsg);
						//log.Print(fieldDataMap)
						json.NewEncoder(rw).Encode(fieldDataList);

					}
				}

			} else {
				sendError(rw, "unknown msg id in url - " + reqUri);
				return;
			}

		} else {
			sendError(rw, "unknown spec id in url - " + reqUri);
		}

	});


}

func ToJsonList(parsedMsg *spec.ParsedMsg) []ui_data.JsonFieldDataRep {


	fieldDataList := make([]ui_data.JsonFieldDataRep, 0, 10);
	for id, fieldData := range (parsedMsg.FieldDataMap) {
		//log.Print(fieldData.Field.Name, fieldData.Value())
		dataRep := ui_data.JsonFieldDataRep{Id:id, Value:fieldData.Field.ValueToString(fieldData.Data)};
		if (fieldData.Field.FieldInfo.Type == spec.BITMAP) {
			dataRep.Value = fieldData.Bitmap.BinaryString();

		}

		fieldDataList = append(fieldDataList, dataRep);
	}

	return fieldDataList;
}

