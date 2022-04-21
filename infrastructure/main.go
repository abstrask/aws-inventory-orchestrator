package main

import (
	"log"
	"os"

	"github.com/dfds/aws-inventory-orchestrator/infrastructure/aws"
	"github.com/dfds/aws-inventory-orchestrator/infrastructure/util"
)

const (
	billingAwsProfile string = "billing-admin"
	runnerAwsProfile  string = "inventory-runner-admin"
)

func main() {

	cliBucketKeyPrefix := "aws/iam/inventory-role/"

	// Parse arguments (IAM role names)
	if len(os.Args) != 3 {
		log.Fatal("Need two arguments: Orchestrator role name and Runner role name")
	}
	orchRole := os.Args[1]
	runnerRole := os.Args[2]

	// Retrieve environment variables
	envMap := util.MapEnvs([]string{"BILLING_ACCOUNT_ID", "INVENTORY_BUCKET_NAME", "INVENTORY_ROLE", "OIDC_PROVIDER", "RUNNER_ACCOUNT_ID", "CLI_BUCKET_NAME"})

	emptyTemplateTokens := util.TemplateTokens{}
	fallbackTemplateTokens := util.TemplateTokens{AccountId: envMap["RUNNER_ACCOUNT_ID"]}

	/* ORCHESTRATOR ROLES */

	// Create PROD inventory orchestrator role
	orchTrustTokens := util.TemplateTokens{AccountId: envMap["BILLING_ACCOUNT_ID"], Oidc: envMap["OIDC_PROVIDER"], ServiceAccount: "aws-inventory-orchestrator-sa"}
	orchTrustDoc := util.ParseTemplateFile("./policies/orchestrator_trust.json", orchTrustTokens)
	orchPolicyDoc := util.ParseTemplateFile("./policies/orchestrator_policy.json", emptyTemplateTokens)
	aws.IamCreateRole(billingAwsProfile, orchRole, "", orchTrustDoc, orchPolicyDoc, 3600)

	// Create TEST inventory orchestrator role (with no trust policy, as it will be managed manually)
	orchRoleNameTest := orchRole + "-Test"
	orchTrustDocTest := util.ParseTemplateFile("./policies/fallback_trust.json", fallbackTemplateTokens)
	aws.IamCreateRole(billingAwsProfile, orchRoleNameTest, "", orchTrustDocTest, orchPolicyDoc, 3600)

	/* RUNNER ROLES */

	// Create PROD inventory runner role
	runnerTrustTokens := util.TemplateTokens{AccountId: envMap["RUNNER_ACCOUNT_ID"], Oidc: envMap["OIDC_PROVIDER"], ServiceAccount: "aws-inventory-runner-sa"}
	runnerTrustDoc := util.ParseTemplateFile("./policies/runner_trust.json", runnerTrustTokens)
	runnerPolicyTokens := util.TemplateTokens{BucketName: envMap["INVENTORY_BUCKET_NAME"], InventoryRole: envMap["INVENTORY_ROLE"]}
	runnerPolicyDoc := util.ParseTemplateFile("./policies/runner_policy.json", runnerPolicyTokens)
	aws.IamCreateRole(runnerAwsProfile, runnerRole, "", runnerTrustDoc, runnerPolicyDoc, 3600)

	// Create TEST inventory runner role (with no trust policy, as it will be managed manually)
	runnerRoleNameTest := runnerRole + "-Test"
	runnerTrustDocTest := util.ParseTemplateFile("./policies/fallback_trust.json", fallbackTemplateTokens)
	aws.IamCreateRole(runnerAwsProfile, runnerRoleNameTest, "", runnerTrustDocTest, runnerPolicyDoc, 3600)

	// Create inventory runner role
	aws.S3CreateBucket(runnerAwsProfile, envMap["INVENTORY_BUCKET_NAME"])

	/* UPLOAD CLI FILES TO S3 BUCKET */
	inventoryTrustTokens := util.TemplateTokens{AccountId: envMap["RUNNER_ACCOUNT_ID"]}
	inventoryTrustDoc := util.ParseTemplateFile("./policies/inventory_trust.json", inventoryTrustTokens)
	inventoryPolicyDoc := util.ParseTemplateFile("./policies/inventory_policy.json", emptyTemplateTokens)
	inventoryPropertiesDoc := util.ParseTemplateFile("./policies/inventory_properties.json", emptyTemplateTokens)
	aws.UploadStringToS3File(runnerAwsProfile, envMap["CLI_BUCKET_NAME"], cliBucketKeyPrefix+"trust.json", inventoryTrustDoc)
	aws.UploadStringToS3File(runnerAwsProfile, envMap["CLI_BUCKET_NAME"], cliBucketKeyPrefix+"policy.json", inventoryPolicyDoc)
	aws.UploadStringToS3File(runnerAwsProfile, envMap["CLI_BUCKET_NAME"], cliBucketKeyPrefix+"properties.json", inventoryPropertiesDoc)
}
