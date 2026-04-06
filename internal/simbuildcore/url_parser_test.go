package simbuildcore

import (
	"net/url"
	"testing"
)

func TestParseAndValidateSimURL_Valid(t *testing.T) {
	t.Setenv(simURLHostsEnv, "simulator.example")
	raw := "https://simulator.example/sim/#skills=Maximum%20Might%20Lv3%2CScorcher%20I&w=LV3-3-3%20Slot%20Weapon&wgs=Lord%27s%20Soul&wss=Gore%20Magala%27s%20Tyranny"

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
	t.Setenv(simURLHostsEnv, "simulator.example")
	raw := "https://simulator.example/sim/#skills=Attack%20Boost%20Lv2&ws=Critical%20Boost%20Lv3%2CMaster's%20Touch%20Lv1"
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
	t.Setenv(simURLHostsEnv, "simulator.example")
	valid, _ := url.Parse("https://simulator.example/sim/#skills=Attack%20Boost%20Lv2")
	if !isSupportedSimURL(valid) {
		t.Fatalf("expected URL to be supported")
	}

	invalid, _ := url.Parse("https://example.com/sim/#skills=Attack%20Boost%20Lv2")
	if isSupportedSimURL(invalid) {
		t.Fatalf("expected URL to be rejected for unsupported host")
	}
}

func TestParseAndValidateSimURL_WithoutScheme(t *testing.T) {
	t.Setenv(simURLHostsEnv, "simulator.example")
	raw := "simulator.example/sim/#skills=Attack%20Boost%20Lv2"
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

func TestIsSupportedSimURL_RejectsWhenAllowListUnset(t *testing.T) {
	t.Setenv(simURLHostsEnv, "")
	valid, _ := url.Parse("https://example-any-host.com/sim/#skills=Attack%20Boost%20Lv2")
	if isSupportedSimURL(valid) {
		t.Fatalf("expected URL to be rejected when host allow-list is not configured")
	}
}
