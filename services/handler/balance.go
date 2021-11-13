package handler

//
//const balanceMinPeriodSeconds = 90
//
//type Balance struct {
//	balance     float64
//	outstanding float64
//	last        time.Time
//	mux         sync.RWMutex
//	check       sync.Mutex
//}
//
//func (h *Handler) BalanceJob() {
//	//for {
//	//	h.GetBalance()
//	//	time.Sleep(time.Minute)
//	//}
//}
//func (b *Balance) CheckBegin() {
//	b.check.Lock()
//}
//func (b *Balance) CheckDone() {
//	b.check.Unlock()
//}
//
//func (b *Balance) Get() (float64, bool) {
//	b.mux.RLock()
//	defer b.mux.RUnlock()
//	var isFresh bool
//	if time.Since(b.last).Seconds() < balanceMinPeriodSeconds {
//		isFresh = true
//	}
//	return b.balance, isFresh
//}
//
//func (b *Balance) Set(bal float64) {
//	b.mux.Lock()
//	defer b.mux.Unlock()
//	b.last = time.Now()
//	b.balance = bal
//}
//func (b *Balance) SetOutstanding(value float64) {
//	b.mux.Lock()
//	defer b.mux.Unlock()
//	b.outstanding = value
//}
//
//func (b *Balance) Sub(risk float64) {
//	b.mux.Lock()
//	defer b.mux.Unlock()
//	b.balance = b.balance - risk
//}
//
////func (h *Handler) GetBalance() float64 {
////	got, b := h.balance.Get()
////	if !b {
////		go h.CheckBalance(false)
////	}
////	return got
////}
//func (b *Balance) FullBalance() float64 {
//	b.mux.RLock()
//	defer b.mux.RUnlock()
//	return util.TruncateFloat(b.balance+b.outstanding, 2)
//}
//
//func (b *Balance) CalcFillFactor() float64 {
//	return util.TruncateFloat(b.outstanding/b.FullBalance(), 2)
//}
//
//func (h *Handler) CheckBalance(force bool) {
//	//h.log.Info("got check balance run")
//	h.balance.CheckBegin()
//	defer h.balance.CheckDone()
//	//h.log.Info("begin check balance")
//	_, b := h.balance.Get()
//	if b && !force {
//		return
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
//	defer cancel()
//	resp, err := h.client.GetBalance(ctx)
//	if err != nil {
//		h.log.Error(err)
//		return
//	}
//	//h.log.Infow("send balance request", "resp", resp)
//
//	h.balance.Set(resp.AvailableBalance)
//	h.balance.SetOutstanding(resp.OutstandingTransactions)
//
//	//full := util.TruncateFloat(resp.AvailableBalance + resp.OutstandingTransactions, 2)
//	//balanceFillFactor:=util.TruncateFloat(resp.AvailableBalance/full, 3)
//	h.log.Debugw("got_balance", "available", resp.AvailableBalance,
//		"outstanding", resp.OutstandingTransactions,
//		"full", h.balance.FullBalance(),
//		"fill_factor", h.balance.CalcFillFactor(),
//		"bet_status", BettingStatus)
//
//	//err = h.store.SaveBalance(ctx, resp, h.AccId())
//	//if err != nil {
//	//	h.log.Info(err)
//	//	return
//	//}
//	//h.log.Info("saved balance")
//}
