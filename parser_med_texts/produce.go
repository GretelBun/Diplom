package main

import (
	"encoding/json"
	"strings"
	"log"
)

type Dialog struct {
	Patient string `json: "patient"`
	Doctor string `json: "doctor"`
}

type ProducedClass struct {
	Name string `json: "name"`
	PatientOnly []Dialog `json: "patient_only"`
	DoctorOnly []Dialog `json: "doctor_only"`
	Both []Dialog `json: "both"`
}

type ProducedClasses map[string]*ProducedClass

type JsonDialog struct {
	Utterances []string `json: "utterances"`
}

func GetJsonDialogs (data []byte) (result []JsonDialog) {
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func splitToDialog(quotes []string) Dialog {
	var patient string
	var doctor string
	for _, quote := range quotes {
		if strings.HasPrefix(quote, "patient") {
			patient = quote
		} else if strings.HasPrefix(quote, "doctor") {
			doctor = quote
		}
	}
	return Dialog{
		Patient: patient,
		Doctor: doctor,
	}
}

func GetDialogs (jsonDialogs []JsonDialog) []Dialog {
	dialogs := make([]Dialog, 0, len(jsonDialogs))

	cnt := 0
	for _, jsonDialog := range jsonDialogs {
		if len(jsonDialog.Utterances) != 2 {
			cnt += 1
			log.Println(cnt)
			continue
		}

		dialogs = append(dialogs, splitToDialog(jsonDialog.Utterances))
	}

	return dialogs
}

func toSlice (classes ProducedClasses) []ProducedClass {
	result := make([]ProducedClass, 0, len(classes))
	for _, val := range classes {
		result = append(result, *val)
	}

	return result
}

func (classesMap ProducedClasses) produce(dialog Dialog, classes []Class) {
	for _, class := range classes {
		var inPatient bool
		var inDoctor bool

		for _, syn := range class.Synonyms {
			inPatient = inPatient || strings.Contains(strings.ToLower(dialog.Patient), syn)
			inDoctor = inDoctor || strings.Contains(strings.ToLower(dialog.Doctor), syn)
			if inPatient && inDoctor {
				break
			}
		}

		if !(inPatient || inDoctor) {
			continue
		}

		if _, ok := classesMap[class.Name]; !ok {
			classesMap[class.Name] = &ProducedClass{Name: class.Name,}
		}

		producedClass := classesMap[class.Name]

		if inPatient && inDoctor {
			producedClass.Both = append(producedClass.Both, dialog)
		} else if inPatient {
			producedClass.PatientOnly = append(producedClass.PatientOnly, dialog)
		} else if inDoctor {
			producedClass.DoctorOnly = append(producedClass.DoctorOnly, dialog)
		}
	}
}

func produceClasses(dialogs []Dialog, classes []Class) ProducedClasses {
	producedClasses := make(ProducedClasses)
	for _, dialog := range dialogs {
		producedClasses.produce(dialog, classes)
	}

	return producedClasses
}

func ProduceClasses(dialogs []Dialog, classes []Class) []ProducedClass {
	return toSlice(produceClasses(dialogs, classes))
}

func GetData(classes []ProducedClass) []byte {
	data, err := json.Marshal(classes)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
