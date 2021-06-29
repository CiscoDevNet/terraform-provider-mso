package mso

import (
	"fmt"
	"log"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateAnp_Basic(t *testing.T) {
	var s SchemaTemplateAnpTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaTemplateAnpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic("test1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
					testAccCheckMsoSchemaTemplateAnpAttributes("test1", &s),
				),
			},
			{
				//Config:            testSchemaTemplateAnpConfig("anp123"),
				ResourceName:      "mso_schema_template_anp.anp1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMSOSchemaTemplateAnp_Basic2(t *testing.T) {
	var s SchemaTemplateAnpTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaTemplateAnpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic2("test1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
					testAccCheckMsoSchemaTemplateAnpAttributes("test1", &s),
				),
			},
			{
				//Config:            testSchemaTemplateAnpConfig("anp123"),
				ResourceName:      "mso_schema_template_anp.anp1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// func TestAccMSOSchemaTemplateAnp_ImportBasic(t *testing.T) {
// 	var s SchemaTemplateAnpTest
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckMsoSchemaTemplateAnpDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic("anp123"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
// 				),
// 			},
// 			{
// 				ResourceName:      "mso_schema_template_anp.anp1",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// func testAccCheckMsoSchemaTemplateAnpConfig_importbasic(name string) string {
// 	return fmt.Sprintf(`

// 	resource "mso_schema" "schema12" {
// 		name = "Schema2"
// 		template_name = "Template1"
// 		tenant_id = "5fb5fed8520000452a9e8911"

// 	  }
// 	  resource "mso_schema_template_anp" "anp12" {
// 		schema_id=mso_schema.schema12.id
// 		template= "Template1"
// 		name = "anp123"
// 		display_name="%s"
// 	  }
// 	`, name)
// }

// func TestAccMSOSchemaTemplateAnp_ImportBasic(t *testing.T) {
// 	var s SchemaTemplateAnpTest
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckMsoSchemaTemplateAnpDestroy,
// 		Steps: []resource.TestStep{
// 			{

// 				Check: resource.ComposeTestCheckFunc(
// 					testSchemaTemplateAnpConfig("anp123", &s),
// 				),
// 			},
// 			// {
// 			// 	Config:            testSchemaTemplateAnpConfig("anp123"),
// 			// 	ResourceName:      "mso_schema_template_anp.anp1",
// 			// 	ImportState:       true,
// 			// 	ImportStateVerify: true,
// 			// },
// 		},
// 	})
// }

// func TestAccMSOSchemaTemplateAnp_Disappears(t *testing.T) {
// 	var s SchemaTemplateAnpTest

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckMsoSchemaTemplateAnpDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic("test1"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
// 					testAccCheckMsoSchemaTemplateAnpAttributes("test1", &s),
// 				),
// 			},
// 			{
// 				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic("test2"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
// 					testAccCheckMsoSchemaTemplateAnpAttributes("test2", &s),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccMSOSchemaTemplateAnp_Update(t *testing.T) {
	var s SchemaTemplateAnpTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaTemplateAnpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic("test1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
					testAccCheckMsoSchemaTemplateAnpAttributes("test1", &s),
				),
			},
			{
				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic("test2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
					testAccCheckMsoSchemaTemplateAnpAttributes("test2", &s),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateAnp_Update2(t *testing.T) {
	var s SchemaTemplateAnpTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaTemplateAnpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic2("test1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
					testAccCheckMsoSchemaTemplateAnpAttributes("test1", &s),
				),
			},
			{
				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic2("test2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
					testAccCheckMsoSchemaTemplateAnpAttributes("test2", &s),
				),
			},
		},
	})
}

// func TestAccMSOSchemaTemplateAnp_Update2(t *testing.T) {
// 	var s SchemaTemplateAnpTest

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckMsoSchemaTemplateAnpDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic("test1"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
// 					testAccCheckMsoSchemaTemplateAnpAttributes("test1", &s),
// 				),
// 			},
// 			{
// 				Config: testAccCheckMsoSchemaTemplateAnpConfig_basic("test2"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMsoSchemaTemplateAnpExists("mso_schema.schema1", "mso_schema_template_anp.anp1", &s),
// 					testAccCheckMsoSchemaTemplateAnpAttributes("test2", &s),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckMsoSchemaTemplateAnpConfig_basic(name string) string {
	return fmt.Sprintf(`

	resource "mso_schema" "schema1" {
		name = "Schema2"
		template_name = "Template1"
		tenant_id = "5fb5fed8520000452a9e8911"
		
	  }
	  resource "mso_schema_template_anp" "anp1" {
		schema_id=mso_schema.schema1.id
		template= "Template1"
		name = "anp123"
		display_name="%s"
	  }
	`, name)
}

func testAccCheckMsoSchemaTemplateAnpConfig_basic2(name string) string {
	return fmt.Sprintf(`

	resource "mso_schema" "schema1" {
		name = "Schema2"
		template_name = "Template1"
		tenant_id = "5fb5fed8520000452a9e8911"
		
	  }
	  resource "mso_schema_template_anp" "anp1" {
		schema_id=mso_schema.schema1.id
		template= "Template2"
		name = "anp123"
		display_name="%s"
	  }
	`, name)
}

func testAccCheckMsoSchemaTemplateAnpConfig_basic3(name string) string {
	return fmt.Sprintf(`

	resource "mso_schema" "schema1" {
		name = "Schema2"
		template_name = "Template1"
		tenant_id = "5fb5fed8520000452a9e8911"
		
	  }
	  resource "mso_schema_template_anp" "anp1" {
		schema_id=mso_schema.schema1.id
		template= "Template2"
		name = "anp123"
		display_name="%s"
	  }
	`, name)
}

// func testSchemaTemplateAnpConfig(name string, stv *SchemaTemplateAnpTest) resource.TestCheckFunc {
// 	client := testAccProvider.Meta().(*client.Client)
// 	cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", "60c1d2255100002f1b4a34ce"))
// 	if err != nil {
// 		return nil
// 	}
// 	// d := SchemaTemplateAnpTest{}
// 	// d.SchemaId := "60c1d2255100002f1b4a34ce"
// 	count, err := cont.ArrayCount("templates")
// 	if err != nil {
// 		return nil
// 	}

// 	for i := 0; i < count; i++ {

// 		tempCont, err := cont.ArrayElement(i, "templates")
// 		if err != nil {
// 			return nil
// 		}
// 		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

// 		if currentTemplateName == "Template1" {
// 			anpCount, err := tempCont.ArrayCount("anps")

// 			if err != nil {
// 				return nil
// 			}
// 			for j := 0; j < anpCount; j++ {
// 				anpCont, err := tempCont.ArrayElement(j, "anps")

// 				if err != nil {
// 					return nil
// 				}
// 				currentAnpName := models.StripQuotes(anpCont.S("name").String())
// 				if currentAnpName == name {
// 					if anpCont.Exists("displayName") {
// 						return nil
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

func testAccCheckMsoSchemaTemplateAnpExists(schemaName string, schemaTemplateAnpName string, stv *SchemaTemplateAnpTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[schemaName]
		rs2, err2 := s.RootModule().Resources[schemaTemplateAnpName]

		if !err1 {
			return fmt.Errorf("Schema %s not found", schemaName)
		}
		if !err2 {
			return fmt.Errorf("Schema Template Anp record %s not found", schemaTemplateAnpName)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No Schema Template Anp id was set")
		}

		client := testAccProvider.Meta().(*client.Client)
		con, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)
		log.Print("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")

		if err != nil {
			return err
		}

		stvt := SchemaTemplateAnpTest{}
		stvt.SchemaId = rs1.Primary.ID

		count, err := con.ArrayCount("templates")
		if err != nil {
			return err
		}

		for i := 0; i < count; i++ {
			tempCont, err := con.ArrayElement(i, "templates")
			stvt.Template = models.StripQuotes(tempCont.S("name").String())
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				if anpCont.Exists("name") {
					stvt.Name = models.StripQuotes(anpCont.S("name").String())

				}

				if anpCont.Exists("displayName") {
					stvt.DisplayName = models.StripQuotes(anpCont.S("displayName").String())
				}

			}
		}

		stvc := &stvt
		log.Printf(fmt.Sprint(stvt.DisplayName))
		*stv = *stvc
		return nil
	}
}

// func testAccCheckMsoSchemaTemplateAnpExists(schemaName string, schemaTemplateAnpName string, stv *SchemaTemplateAnpTest) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs1, err1 := s.RootModule().Resources[schemaName]
// 		rs2, err2 := s.RootModule().Resources[schemaTemplateAnpName]

// 		if !err1 {
// 			return fmt.Errorf("Schema %s not found", schemaName)
// 		}
// 		if !err2 {
// 			return fmt.Errorf("Schema Template Anp record %s not found", schemaTemplateAnpName)
// 		}

// 		if rs1.Primary.ID == "" {
// 			return fmt.Errorf("No Schema id was set")
// 		}
// 		if rs2.Primary.ID == "" {
// 			return fmt.Errorf("No Schema Template Anp id was set")
// 		}

// 		client := testAccProvider.Meta().(*client.Client)
// 		con, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)

// 		if err != nil {
// 			return err
// 		}

// 		stvt := SchemaTemplateAnpTest{}
// 		stvt.SchemaId = rs1.Primary.ID

// 		count, err := con.ArrayCount("templates")
// 		if err != nil {
// 			return err
// 		}

// 		for i := 0; i < count; i++ {
// 			tempCont, err := con.ArrayElement(i, "templates")
// 			stvt.Template = models.StripQuotes(tempCont.S("name").String())
// 			anpCount, err := tempCont.ArrayCount("anps")
// 			if err != nil {
// 				return fmt.Errorf("No Anp found")
// 			}
// 			for j := 0; j < anpCount; j++ {
// 				anpCont, err := tempCont.ArrayElement(j, "anps")
// 				if err != nil {
// 					return err
// 				}
// 				if anpCont.Exists("name") {
// 					stvt.Name = models.StripQuotes(anpCont.S("name").String())

// 				}

// 				if anpCont.Exists("displayName") {
// 					stvt.DisplayName = models.StripQuotes(anpCont.S("displayName").String())
// 				}

// 			}
// 		}

// 		stvc := &stvt
// 		log.Printf("^^^^^^^^^^^^^^^^^^^")
// 		log.Printf(fmt.Sprintf(stvt.Name))

// 		*stv = *stvc
// 		return nil
// 	}
// }

func testAccCheckMsoSchemaTemplateAnpDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	rs1, err1 := s.RootModule().Resources["mso_schema.schema1"]

	if !err1 {
		return fmt.Errorf("Schema %s not found", "mso_schema.schema1")
	}

	schemaid := rs1.Primary.ID
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_anp" {
			con, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaid))
			if err != nil {
				return nil
			} else {
				count, err := con.ArrayCount("templates")
				if err != nil {
					return fmt.Errorf("No Template found")
				}
				for i := 0; i < count; i++ {
					tempCont, err := con.ArrayElement(i, "templates")
					if err != nil {
						return fmt.Errorf("No template exists")
					}
					anpCount, err := tempCont.ArrayCount("anps")
					if err != nil {
						return fmt.Errorf("No Anp found")
					}
					for j := 0; j < anpCount; j++ {
						anpCont, err := tempCont.ArrayElement(j, "anps")
						if err != nil {
							return err
						}
						name := models.StripQuotes(anpCont.S("name").String())

						if rs.Primary.ID == name {
							return fmt.Errorf("Schema Template Anp record still exists")

						}

					}
				}
			}
		}
	}
	return nil
}

func testAccCheckMsoSchemaTemplateAnpAttributes(name string, stv *SchemaTemplateAnpTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name != stv.DisplayName {
			return fmt.Errorf("Bad Schema Template Anp Name %s", stv.DisplayName)
		}
		return nil
	}
}

type SchemaTemplateAnpTest struct {
	Id          string `json:",omitempty"`
	SchemaId    string `json:",omitempty"`
	Template    string `json:",omitempty"`
	Name        string `json:",omitempty"`
	DisplayName string `json:",omitempty"`
}
