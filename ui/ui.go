// SPDX-License-Identifier: Unlicense OR MIT

package ui

// A Gio program that demonstrates Gio widgets. See https://gioui.org for more information.

import (
	"flag"
	"fmt"
	"log"
	"math"
	"orderbook"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gofrs/uuid"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/font/gofont"
)

var screenshot = flag.String("screenshot", "", "save a screenshot to a file and exit")

type scaledConfig struct {
	Scale float32
}

func Run(ob *orderbook.OrderBook) error {
	flag.Parse()
	editor.SetText(longText)
	ic, err := material.NewIcon(icons.ContentAdd)
	if err != nil {
		return err
	}
	icon = ic
	gofont.Register()

	go func() {
		w := app.NewWindow()
		err := loop(w, ob)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	app.Main()
	return nil
}

func loop(w *app.Window, ob *orderbook.OrderBook) error {
	th := material.NewTheme()
	gtx := layout.NewContext(w.Queue())
	t := time.NewTicker(1 * time.Millisecond)
	go func() {
		for {
			select {
			case <-t.C:
				w.Invalidate()
			}
		}
	}()
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)
			orderbookUI(gtx, th, ob)
			e.Frame(gtx.Ops)
		}
	}
}

var (
	editor       = new(widget.Editor)
	amountEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	priceEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	button            = new(widget.Button)
	greenButton       = new(widget.Button)
	iconButton        = new(widget.Button)
	radioButtonsGroup = new(widget.Enum)
	list              = &layout.List{
		Axis: layout.Vertical,
	}
	buyList = &layout.List{
		Axis: layout.Vertical,
	}
	sellList = &layout.List{
		Axis: layout.Vertical,
	}
	green    = true
	icon     *material.Icon
	checkbox = new(widget.CheckBox)
)

func renderOrder(blank bool, order *orderbook.Order) string {
	tmpl := "\t\t%03d @ %03d USD"
	if blank {
		return "\t\t" + "--- @ --- USD"
	}
	return fmt.Sprintf(tmpl, order.Amount, order.Price)
}
func orderbookUI(gtx *layout.Context, th *material.Theme, ob *orderbook.OrderBook) {

	widgets := []func(){
		func() {
			th.H3("Orderbook").Layout(gtx)
		},
		func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func() {
					th.RadioButton("buy", "Buy").Layout(gtx, radioButtonsGroup)
				}),
				layout.Rigid(func() {
					th.RadioButton("sell", "Sell").Layout(gtx, radioButtonsGroup)
				}),
				layout.Rigid(func() {
					e := th.Editor("Amount")
					e.Font.Style = text.Italic
					e.Layout(gtx, amountEditor)
				}),
				layout.Rigid(func() {
					e := th.Editor("Price")
					e.Font.Style = text.Italic
					e.Layout(gtx, priceEditor)
				}),
				layout.Rigid(func() {
					layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
						for button.Clicked(gtx) {

							side := radioButtonsGroup.Value(gtx)
							if side == "" {
								fmt.Println("select side first")
								continue
							}
							amount, err := strconv.Atoi(amountEditor.Text())
							if err != nil {
								fmt.Println(err)
								continue
							}
							price, err := strconv.Atoi(priceEditor.Text())
							if err != nil {
								fmt.Println(err)
								continue
							}
							OrderSide := orderbook.SELL
							if side == "buy" {
								OrderSide = orderbook.BUY
							}
							ob.Process(&orderbook.Order{
								ID:     uuid.Must(uuid.NewV4()).String(),
								Amount: uint64(amount),
								Price:  uint64(price),
								Side:   OrderSide,
							})
						}
						th.Button("Click me!").Layout(gtx, button)
					})
				}),
			)
		},
		func() {
			gtx.Constraints.Height.Min = gtx.Px(unit.Dp(200))
			orders := []string{}
			for i := 0; i < 10; i++ {
				n := len(ob.SellOrders)
				j := 9 - i
				if j > n-1 {
					orders = append(orders, renderOrder(true, &orderbook.Order{}))
					continue
				}
				orders = append(orders, renderOrder(false, ob.SellOrders[j]))
			}
			orders = append(orders, "\n")
			for i := 0; i < 10; i++ {
				n := len(ob.BuyOrders)
				if i > n-1 {
					orders = append(orders, renderOrder(true, &orderbook.Order{}))
					continue
				}
				orders = append(orders, renderOrder(false, ob.BuyOrders[i]))
			}
			th.Body1(strings.Join(orders, "\n")).Layout(gtx)
		},
	}
	list.Layout(gtx, len(widgets), func(i int) {
		layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}

func (s *scaledConfig) Now() time.Time {
	return time.Now()
}

func (s *scaledConfig) Px(v unit.Value) int {
	scale := s.Scale
	if v.U == unit.UnitPx {
		scale = 1
	}
	return int(math.Round(float64(scale * v.V)))
}

const longText = `1. I learned from my grandfather, Verus, to use good manners, and to
put restraint on anger. 2. In the famous memory of my father I had a
pattern of modesty and manliness. 3. Of my mother I learned to be
pious and generous; to keep myself not only from evil deeds, but even
from evil thoughts; and to live with a simplicity which is far from
customary among the rich. 4. I owe it to my great-grandfather that I
did not attend public lectures and discussions, but had good and able
teachers at home; and I owe him also the knowledge that for things of
this nature a man should count no expense too great.

5. My tutor taught me not to favour either green or blue at the
chariot races, nor, in the contests of gladiators, to be a supporter
either of light or heavy armed. He taught me also to endure labour;
not to need many things; to serve myself without troubling others; not
to intermeddle in the affairs of others, and not easily to listen to
slanders against them.

6. Of Diognetus I had the lesson not to busy myself about vain things;
not to credit the great professions of such as pretend to work
wonders, or of sorcerers about their charms, and their expelling of
Demons and the like; not to keep quails (for fighting or divination),
nor to run after such things; to suffer freedom of speech in others,
and to apply myself heartily to philosophy. Him also I must thank for
my hearing first Bacchius, then Tandasis and Marcianus; that I wrote
dialogues in my youth, and took a liking to the philosopher's pallet
and skins, and to the other things which, by the Grecian discipline,
belong to that profession.

7. To Rusticus I owe my first apprehensions that my nature needed
reform and cure; and that I did not fall into the ambition of the
common Sophists, either by composing speculative writings or by
declaiming harangues of exhortation in public; further, that I never
strove to be admired by ostentation of great patience in an ascetic
life, or by display of activity and application; that I gave over the
study of rhetoric, poetry, and the graces of language; and that I did
not pace my house in my senatorial robes, or practise any similar
affectation. I observed also the simplicity of style in his letters,
particularly in that which he wrote to my mother from Sinuessa. I
learned from him to be easily appeased, and to be readily reconciled
with those who had displeased me or given cause of offence, so soon as
they inclined to make their peace; to read with care; not to rest
satisfied with a slight and superficial knowledge; nor quickly to
assent to great talkers. I have him to thank that I met with the
discourses of Epictetus, which he furnished me from his own library.

8. From Apollonius I learned true liberty, and tenacity of purpose; to
regard nothing else, even in the smallest degree, but reason always;
and always to remain unaltered in the agonies of pain, in the losses
of children, or in long diseases. He afforded me a living example of
how the same man can, upon occasion, be most yielding and most
inflexible. He was patient in exposition; and, as might well be seen,
esteemed his fine skill and ability in teaching others the principles
of philosophy as the least of his endowments. It was from him that I
learned how to receive from friends what are thought favours without
seeming humbled by the giver or insensible to the gift.`
