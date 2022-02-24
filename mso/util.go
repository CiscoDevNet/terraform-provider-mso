package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const version = 1

func toStringList(configured interface{}) []string {
	vs := make([]string, 0, 1)
	val, ok := configured.(string)
	if ok && val != "" {
		vs = append(vs, val)
	}
	return vs
}

func makeTestVariable(s string) string {
	return fmt.Sprintf("acctest_%s", s)
}

func epgRefValidation() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		res, err := regexp.MatchString(`^\/schemas\/(.)+\/templates\/(.)+\/anps\/(.)+\/epgs\/(.)+$`, v)
		if !res {
			log.Printf("err: %v\n", err)
			es = append(es, fmt.Errorf("invalid epg reference:expected format /schema/{schema_id}/templates/{template_name}/anps/{anp_name}/epgs/{epg_name}"))
		}
		return
	}
}

func externalEpgRefValidation() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		res, err := regexp.MatchString(`^\/schemas\/(.)+\/templates\/(.)+\/externalEpgs\/(.)+$`, v)
		if !res {
			log.Printf("err: %v\n", err)
			es = append(es, fmt.Errorf("invalid external epg reference:expected format /schema/{schema_id}/template/{template_name}/externalEpgs/{external_epg_name}"))
		}
		return
	}
}
