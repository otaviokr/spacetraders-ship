package web

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "spacetradership"
)

var (
	TradeCycles = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "trade_cycles",
			Help:      "How many times the trade route has been performed",
		},
		[]string{"ship_id"})

	FuelConsumed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "fuel_consumed",
			Help:      "How much fuel has been consummed",
		},
		[]string{"ship_id"})

	MoneySpent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "money_spent",
			Help:      "How many credits the starship used to buy goods",
		},
		[]string{"ship_id", "good"})

	MoneyEarned = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "money_earned",
			Help:      "How many credits the starship received from selling goods",
		},
		[]string{"ship_id", "good"})

	GoodsSold = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "goods_sold",
			Help:      "Products sold",
		},
		[]string{"ship_id", "good", "location"})

	GoodsBought = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "goods_bought",
			Help:      "Products bought",
		},
		[]string{"ship_id", "good", "location"})
)
