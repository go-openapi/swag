// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package mangling

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestSplitPluralized(t *testing.T) {
	m := NewNameMangler(
		WithAdditionalInitialisms("elb", "cap", "capwd", "wd"),
	)

	s := newSplitter(
		withInitialismsCache(&m.index.initialismsCache),
		withPostSplitInitialismCheck,
	)

	t.Run("should recognize pluralized initialisms", func(t *testing.T) {
		t.Run("with trailing initialism", func(t *testing.T) {
			const plurals = "pluralized initialism IDs"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 3)

			assert.Equal(t, "PluralizedInitialismIDs", m.ToGoName(plurals))
			assert.Equal(t, "pluralized_initialism_ids", m.ToFileName(plurals))
		})

		t.Run("with initialism trailed by capital", func(t *testing.T) {
			const plurals = "pluralized initialism IDsX"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 4)

			assert.Equal(t, "PluralizedInitialismIDsX", m.ToGoName(plurals))
			assert.Equal(t, "pluralized_initialism_ids_x", m.ToFileName(plurals))
		})

		t.Run("with middle initialism", func(t *testing.T) {
			const plurals = "pluralized IDs initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 3)

			assert.Equal(t, "PluralizedIDsInitialism", m.ToGoName(plurals))
			assert.Equal(t, "pluralized_ids_initialism", m.ToFileName(plurals))
		})

		t.Run("with upper-cased pluralized initialism", func(t *testing.T) {
			const plurals = "pluralized IDS initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 4)

			assert.Equal(t, "PluralizedIDSInitialism", m.ToGoName(plurals))
			assert.Equal(t, "pluralized_id_s_initialism", m.ToFileName(plurals))
		})

		t.Run("with leading initialism", func(t *testing.T) {
			const plurals = "IDs pluralized initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 3)

			assert.Equal(t, "IDsPluralizedInitialism", m.ToGoName(plurals))
			assert.Equal(t, "ids_pluralized_initialism", m.ToFileName(plurals))
		})

		t.Run("with added non-default initialisms", func(t *testing.T) {
			const plurals = "pluralized initialism ELBs"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 3)

			assert.Equal(t, "PluralizedInitialismELBs", m.ToGoName(plurals))
			assert.Equal(t, "pluralized_initialism_elbs", m.ToFileName(plurals))
		})
	})

	t.Run("should recognize invariant initialisms", func(t *testing.T) {
		t.Run("with explicit word boundary", func(t *testing.T) {
			const plurals = "pluralized HTTP's initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 4)

			assert.Equal(t, "PluralizedHTTPsInitialism", m.ToGoName(plurals))
			assert.Equal(t, "pluralized_http_s_initialism", m.ToFileName(plurals))
		})

		t.Run("with continued word", func(t *testing.T) {
			t.Run("no initialism (invariant)", func(t *testing.T) {
				const plurals = "pluralized HTTPs is not an initialism"
				lexems := s.split(plurals)
				poolOfLexems.RedeemLexems(lexems)

				require.NotNil(t, lexems)
				require.Len(t, *lexems, 9)
				assert.Equal(t, "PluralizedHTTPsIsNotAnInitialism", m.ToGoName(plurals))
				assert.Equal(t, "pluralizedHTTPsIsNotAnInitialism", m.ToVarName(plurals))
				assert.Equal(t, "pluralized_h_t_t_ps_is_not_an_initialism", m.ToFileName(plurals))
			})

			t.Run("no initialism (pluralizable)", func(t *testing.T) {
				const plurals = "pluralized ELBsis not an initialism"
				lexems := s.split(plurals)
				poolOfLexems.RedeemLexems(lexems)

				require.NotNil(t, lexems)
				require.Len(t, *lexems, 7)
				assert.Equal(t, "PluralizedELBsisNotAnInitialism", m.ToGoName(plurals))
				assert.Equal(t, "pluralizedELBsisNotAnInitialism", m.ToVarName(plurals))
				assert.Equal(t, "pluralized_e_l_bsis_not_an_initialism", m.ToFileName(plurals))
			})

			t.Run("no initialism (no plural)", func(t *testing.T) {
				const plurals = "pluralized ELBx is not an initialism"
				lexems := s.split(plurals)
				poolOfLexems.RedeemLexems(lexems)

				require.NotNil(t, lexems)
				require.Len(t, *lexems, 8)
				assert.Equal(t, "PluralizedELBxIsNotAnInitialism", m.ToGoName(plurals))
				assert.Equal(t, "pluralizedELBxIsNotAnInitialism", m.ToVarName(plurals))
				assert.Equal(t, "pluralized_e_l_bx_is_not_an_initialism", m.ToFileName(plurals))
			})

			t.Run("with initialism trailed by lowercase", func(t *testing.T) {
				const plurals = "pluralized initialism IDsx"
				lexems := s.split(plurals)
				poolOfLexems.RedeemLexems(lexems)

				require.NotNil(t, lexems)
				require.Len(t, *lexems, 4)

				assert.Equal(t, "PluralizedInitialismIDsx", m.ToGoName(plurals))
				assert.Equal(t, "pluralized_initialism_i_dsx", m.ToFileName(plurals))
			})
		})

		t.Run("with proper case match: detect initialism", func(t *testing.T) {
			const plurals = "pluralized HTTPS is an initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 5)
			assert.Equal(t, "PluralizedHTTPSIsAnInitialism", m.ToGoName(plurals))
			assert.Equal(t, "pluralizedHTTPSIsAnInitialism", m.ToVarName(plurals))
			assert.Equal(t, "pluralized_https_is_an_initialism", m.ToFileName(plurals))
		})
	})
}

func TestSplitter(t *testing.T) {
	s := newSplitter(withPostSplitInitialismCheck)

	t.Run("should return an empty slice of lexems", func(t *testing.T) {
		lexems := s.split("")
		poolOfLexems.RedeemLexems(lexems)

		require.NotNil(t, lexems)
		require.Empty(t, lexems)
	})
}
