create or alter view dbo.GetResults as
select top 2200 b.SurebetId,
                b.SideIndex,
                b.BetId,
                b.ApiBetId,
                l.BetStatus ApiBetStatus,
                l.Price     Price,
                l.Risk      Stake,
                l.WinLoss   WinLoss
from Bet b
         left join StraightBetList l on b.ApiBetId = l.BetId
where b.ApiBetId > 0
  and l.BetStatus != 'ACCEPTED'
order by SurebetId desc

