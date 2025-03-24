package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitPluralized(t *testing.T) {
	s := newSplitter(withPostSplitInitialismCheck)

	t.Run("should recognize pluralized initialisms", func(t *testing.T) {
		t.Run("with trailing initialism", func(t *testing.T) {
			const plurals = "pluralized initialism IDs"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 3)

			assert.Equal(t, "PluralizedInitialismIDs", ToGoName(plurals))
			assert.Equal(t, "pluralized_initialism_ids", ToFileName(plurals))
		})

		t.Run("with initialism trailed by capital", func(t *testing.T) {
			const plurals = "pluralized initialism IDsX"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 4)

			assert.Equal(t, "PluralizedInitialismIDsX", ToGoName(plurals))
			assert.Equal(t, "pluralized_initialism_ids_x", ToFileName(plurals))
		})

		t.Run("with middle initialism", func(t *testing.T) {
			const plurals = "pluralized IDs initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 3)

			assert.Equal(t, "PluralizedIDsInitialism", ToGoName(plurals))
			assert.Equal(t, "pluralized_ids_initialism", ToFileName(plurals))
		})

		t.Run("with upper-cased pluralized initialism", func(t *testing.T) {
			const plurals = "pluralized IDS initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 4)

			assert.Equal(t, "PluralizedIDSInitialism", ToGoName(plurals))
			assert.Equal(t, "pluralized_id_s_initialism", ToFileName(plurals))
		})

		t.Run("with leading initialism", func(t *testing.T) {
			const plurals = "IDs pluralized initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 3)

			assert.Equal(t, "IDsPluralizedInitialism", ToGoName(plurals))
			assert.Equal(t, "ids_pluralized_initialism", ToFileName(plurals))
		})

		t.Run("with added non-default initialisms", func(t *testing.T) {
			const plurals = "pluralized initialism ELBs"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 3)

			assert.Equal(t, "PluralizedInitialismELBs", ToGoName(plurals))
			assert.Equal(t, "pluralized_initialism_elbs", ToFileName(plurals))
		})
	})

	t.Run("should recognize invariant initialisms", func(t *testing.T) {
		t.Run("with explicit word boundary", func(t *testing.T) {
			const plurals = "pluralized HTTP's initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 4)

			assert.Equal(t, "PluralizedHTTPsInitialism", ToGoName(plurals))
			assert.Equal(t, "pluralized_http_s_initialism", ToFileName(plurals))
		})

		t.Run("with continued word", func(t *testing.T) {
			t.Run("no initialism (invariant)", func(t *testing.T) {
				const plurals = "pluralized HTTPs is not an initialism"
				lexems := s.split(plurals)
				poolOfLexems.RedeemLexems(lexems)

				require.NotNil(t, lexems)
				require.Len(t, *lexems, 9)
				assert.Equal(t, "PluralizedHTTPsIsNotAnInitialism", ToGoName(plurals))
				assert.Equal(t, "pluralizedHTTPsIsNotAnInitialism", ToVarName(plurals))
				assert.Equal(t, "pluralized_h_t_t_ps_is_not_an_initialism", ToFileName(plurals))
			})

			t.Run("no initialism (pluralizable)", func(t *testing.T) {
				const plurals = "pluralized ELBsis not an initialism"
				lexems := s.split(plurals)
				poolOfLexems.RedeemLexems(lexems)

				require.NotNil(t, lexems)
				require.Len(t, *lexems, 7)
				assert.Equal(t, "PluralizedELBsisNotAnInitialism", ToGoName(plurals))
				assert.Equal(t, "pluralizedELBsisNotAnInitialism", ToVarName(plurals))
				assert.Equal(t, "pluralized_e_l_bsis_not_an_initialism", ToFileName(plurals))
			})

			t.Run("no initialism (no plural)", func(t *testing.T) {
				const plurals = "pluralized ELBx is not an initialism"
				lexems := s.split(plurals)
				poolOfLexems.RedeemLexems(lexems)

				require.NotNil(t, lexems)
				require.Len(t, *lexems, 8)
				assert.Equal(t, "PluralizedELBxIsNotAnInitialism", ToGoName(plurals))
				assert.Equal(t, "pluralizedELBxIsNotAnInitialism", ToVarName(plurals))
				assert.Equal(t, "pluralized_e_l_bx_is_not_an_initialism", ToFileName(plurals))
			})

			t.Run("with initialism trailed by lowercase", func(t *testing.T) {
				const plurals = "pluralized initialism IDsx"
				lexems := s.split(plurals)
				poolOfLexems.RedeemLexems(lexems)

				require.NotNil(t, lexems)
				require.Len(t, *lexems, 4)

				assert.Equal(t, "PluralizedInitialismIDsx", ToGoName(plurals))
				assert.Equal(t, "pluralized_initialism_i_dsx", ToFileName(plurals))
			})
		})

		t.Run("with proper case match: detect initialism", func(t *testing.T) {
			const plurals = "pluralized HTTPS is an initialism"
			lexems := s.split(plurals)
			poolOfLexems.RedeemLexems(lexems)

			require.NotNil(t, lexems)
			require.Len(t, *lexems, 5)
			assert.Equal(t, "PluralizedHTTPSIsAnInitialism", ToGoName(plurals))
			assert.Equal(t, "pluralizedHTTPSIsAnInitialism", ToVarName(plurals))
			assert.Equal(t, "pluralized_https_is_an_initialism", ToFileName(plurals))
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

func TestIsEqualFoldIgnoreSpace(t *testing.T) {
	t.Run("should find equal", func(t *testing.T) {
		require.True(t, isEqualFoldIgnoreSpace([]rune(""), ""))
		require.True(t, isEqualFoldIgnoreSpace([]rune(""), "  "))

		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), " a"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), "a "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), " a "))

		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), "\ta\t"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), "a"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), "\u00A0a\u00A0"))

		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " ab"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), "ab "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " ab "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " ab "))

		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), "AB "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB "))

		require.True(t, isEqualFoldIgnoreSpace([]rune("À"), " à"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("À"), "à "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("À"), " à "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("À"), " à "))
	})

	t.Run("should find different", func(t *testing.T) {
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))

		require.False(t, isEqualFoldIgnoreSpace([]rune("AB"), " A B "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("AB"), " a b "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB \u00A0\u00A0x"))
		require.False(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB \u00A0\u00A0é"))

		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))

		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " à"))
		require.False(t, isEqualFoldIgnoreSpace([]rune("À"), " bà"))
		require.False(t, isEqualFoldIgnoreSpace([]rune("À"), "àb "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("À"), " a "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("À"), "Á"))
	})
}
