package functors

import (
	"github.com/influx6/flux"
	"github.com/influx6/frpcore/reactive"
)

type (

	//ObserveFunction represent a listen function
	ObserveFunction func(interface{})

	//GenericFunction represents a generic function type
	GenericFunction func(interface{}) interface{}

	//ImmuteFunction represents a generic function type
	ImmuteFunction func(interface{}) reactive.Immutable

	//Functor represent the base function type
	Functor flux.Stacks
)

//Transform transforms a function into a functor
func Transform(fx GenericFunction) Functor {
	return Functor(flux.NewStack(flux.WrapStackable(fx)))
}

//KindTransform transforms a function into a functor that accepts on a specific kind of value
func KindTransform(ktype interface{}, fx GenericFunction) Functor {
	kind := reactive.GetKind(ktype)
	return Functor(flux.NewStack(flux.WrapStackable(func(bx interface{}) interface{} {
		if reactive.AcceptableKind(kind, bx) {
			return fx(bx)
		}
		return nil
	})))
}

//ImmuteTransform transforms all supplied values into reactive.Immutable that watches for a specific kind
func ImmuteTransform(ktype interface{}, fx ImmuteFunction) Functor {
	return KindTransform(ktype, func(bx interface{}) interface{} {
		imx, ok := bx.(reactive.Immutable)
		if ok {
			return fx(imx.Value())
		}
		return fx(bx)
	})
}

//WrapStacks a list of functors into a list of flux.Stacks
func WrapStacks(n ...Functor) []flux.Stacks {
	var fx []flux.Stacks
	for _, m := range n {
		fx = append(fx, flux.Stacks(m))
	}
	return fx
}

//Lift combines multiple functors
func Lift(n ...Functor) Functor {
	return Functor(flux.FeedAllStacks(false, WrapStacks(n...)...))
}

//Merge combines one value from another functor
func Merge(n ...Functor) Functor {

	return nil
}
