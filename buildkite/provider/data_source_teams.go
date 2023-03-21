package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/tink-ab/terraform-provider-buildkite/buildkite/client"
)

func dataSourceTeams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamsRead,
		Schema: map[string]*schema.Schema{
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"team_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"privacy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      client.TeamPrivacyVisible,
							ValidateFunc: validation.StringInSlice(ValidTeamPrivacy, false),
						},
						"default_member_role": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      client.TeamMemberRoleMember,
							ValidateFunc: validation.StringInSlice(ValidTeamMemberRole, false),
						},
						"is_default_team": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func dataSourceTeamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*client.Client)
	ts, err := client.GetTeams()
	if err != nil {
		return diag.FromErr(err)
	}
	teams := make([]map[string]interface{}, 0)
	for _, t := range *ts {
		teams = append(teams, map[string]interface{}{
			"team_id":             t.Id,
			"uuid":                t.UUID,
			"slug":                t.Slug,
			"name":                t.Name,
			"description":         t.Description,
			"created_at":          t.CreatedAt,
			"privacy":             t.Privacy,
			"is_default_team":     t.IsDefaultTeam,
			"default_member_role": t.DefaultMemberRole,
		})
	}

	if err := d.Set("teams", teams); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
