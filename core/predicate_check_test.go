// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

// TODO
// 3. unit test
// 4. pass over documentation
// 5. change to only use engine context for warp (not all precompiles)

//  var _ precompileconfig.Predicater = (*mockPredicater)(nil)
// type mockPredicater struct {
// 	predicateFunc func(predicateContext *precompileconfig.PredicateContext, b []byte) error
// }

// func (m *mockPredicater) VerifyPredicate(predicateContext *precompileconfig.PredicateContext, b []byte) error { return m.predicateFunc(predicateContext, b) }
// type predicateCheckTest struct {

// }

// func TestCheckPredicate(t *testing.T) {
// 	for name, test := range []predicateCheckTest{

// 	} {
// 		test := test
// 		t.Run(name, func(t *testing.T) {
// 			rules := params.TestChainConfig.AvalancheRules(common.Big0, common.Big0)
// 			rules.PredicatePrecompiles[common.HexToAddress("")]
// 			CheckPredicates()
// 		})
// 	}
// }
