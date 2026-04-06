package simbuildcore

import (
	"net/url"
	"strings"
	"testing"
)

func TestParseAndValidateSimURL_Valid(t *testing.T) {
	raw := "https://mhwilds.wiki-db.com/sim/#skills=Maximum%20Might%20Lv3%2CScorcher%20I&w=LV3-3-3%20Slot%20Weapon&wgs=Lord%27s%20Soul&wss=Gore%20Magala%27s%20Tyranny"

	out, err := parseAndValidateSimURL(raw)
	if err != nil {
		t.Fatalf("expected valid URL, got error: %v", err)
	}

	if len(out.WeaponSetSkills) != 2 {
		t.Fatalf("expected 2 weapon set skills, got %d", len(out.WeaponSetSkills))
	}
	if len(out.Skills) != 4 {
		t.Fatalf("expected 4 parsed skills, got %d", len(out.Skills))
	}

	if out.Skills[0].BaseName != "Maximum Might" || out.Skills[0].RequestedLevel != 3 {
		t.Fatalf("unexpected first skill parsing: %+v", out.Skills[0])
	}
	if out.Skills[1].BaseName != "Scorcher" || out.Skills[1].RequestedLevel != 1 {
		t.Fatalf("unexpected roman-level parsing: %+v", out.Skills[1])
	}
}

func TestParseAndValidateSimURL_ParsesWeaponExtraSkills(t *testing.T) {
	raw := "https://mhwilds.wiki-db.com/sim/#skills=Attack%20Boost%20Lv2&ws=Critical%20Boost%20Lv3%2CMaster's%20Touch%20Lv1"
	out, err := parseAndValidateSimURL(raw)
	if err != nil {
		t.Fatalf("expected valid URL, got error: %v", err)
	}
	if len(out.WeaponExtraSkills) != 2 {
		t.Fatalf("expected two parsed weapon extra skills, got %d", len(out.WeaponExtraSkills))
	}
	if out.WeaponExtraSkills[0].BaseName != "Critical Boost" || out.WeaponExtraSkills[0].RequestedLevel != 3 {
		t.Fatalf("unexpected first weapon extra skill: %+v", out.WeaponExtraSkills[0])
	}
	if out.WeaponExtraSkills[1].BaseName != "Master's Touch" || out.WeaponExtraSkills[1].RequestedLevel != 1 {
		t.Fatalf("unexpected second weapon extra skill: %+v", out.WeaponExtraSkills[1])
	}
}

func TestIsSupportedSimURL(t *testing.T) {
	valid, _ := url.Parse("https://mhwilds.wiki-db.com/sim/#skills=Attack%20Boost%20Lv2")
	if !isSupportedSimURL(valid) {
		t.Fatalf("expected URL to be supported")
	}

	invalid, _ := url.Parse("https://example.com/sim/#skills=Attack%20Boost%20Lv2")
	if isSupportedSimURL(invalid) {
		t.Fatalf("expected URL to be rejected for unsupported host")
	}
}

func TestParseAndValidateSimURL_WithoutScheme(t *testing.T) {
	raw := "mhwilds.wiki-db.com/sim/#skills=Attack%20Boost%20Lv2"
	out, err := parseAndValidateSimURL(raw)
	if err != nil {
		t.Fatalf("expected parser to accept host/path without scheme: %v", err)
	}
	if len(out.Skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(out.Skills))
	}
	if out.Skills[0].BaseName != "Attack Boost" || out.Skills[0].RequestedLevel != 2 {
		t.Fatalf("unexpected parsed skill: %+v", out.Skills[0])
	}
}

func TestParseAndValidateSimURL_UserExampleWithApostrophes(t *testing.T) {
	raw := "https://mhwilds.wiki-db.com/sim/#skills=Maximum%20Might%20Lv3%2CCritical%20Boost%20Lv4&wgs=Lord's%20Soul&wss=Gore%20Magala's%20Tyranny"
	out, err := parseAndValidateSimURL(raw)
	if err != nil {
		t.Fatalf("expected parser to accept user-style URL, got error: %v", err)
	}
	if len(out.WeaponSetSkills) != 2 {
		t.Fatalf("expected 2 parsed weapon set skills, got %d", len(out.WeaponSetSkills))
	}
	if got := out.WeaponSetSkills[0].BaseName; !strings.EqualFold(got, "Lord's Soul") {
		t.Fatalf("unexpected weapon group skill name: %q", got)
	}
	if got := out.WeaponSetSkills[1].BaseName; !strings.EqualFold(got, "Gore Magala's Tyranny") {
		t.Fatalf("unexpected weapon set skill name: %q", got)
	}
}

func TestParseAndValidateSimURL_ExactUserPayloadShape(t *testing.T) {
	raw := "https://mhwilds.wiki-db.com/sim/#skills=Maximum%20Might%20Lv3%2CWeakness%20Exploit%20Lv5%2CAgitator%20Lv5%2CBurst%20Lv1%2CScorcher%20I%2CBlack%20Eclipse%20I%2CGuts%20(Tenacity)%2CAntivirus%20Lv3%2CBinding%20Counter%20I%2COffensive%20Guard%20Lv3%2CWater%20Attack%20Lv1%2CMaster's%20Touch%20Lv1%2CEarplugs%20Lv2%2CCritical%20Eye%20Lv3%2CCritical%20Boost%20Lv4&s=1&e=1&v=10&g=13&w=LV3-3-3%20Slot%20Weapon&ws=&d=0&rf=-100&rw=-100&rt=-100&ri=-100&rd=-100&l=200&wgs=Lord's%20Soul&wss=Gore%20Magala's%20Tyranny"
	out, err := parseAndValidateSimURL(raw)
	if err != nil {
		t.Fatalf("expected parser to accept exact user payload shape, got error: %v", err)
	}
	if len(out.Skills) == 0 {
		t.Fatalf("expected at least one parsed skill")
	}
}
