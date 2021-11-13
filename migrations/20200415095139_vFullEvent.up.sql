create or alter view dbo.vFullEvent as
select S.Id SportId,
       S.Name SportName,
       L.Id LeagueId,
       L.Name LeagueName,
       E.Id EventId,
       E.Starts EventStarts,
       E.Home Home,
       E.Away Away,
       ln.LineId LineId,
       M.Home PriceHome,
       M.Away PriceAway,
       M.Draw PriceDraw
from Sport S
         join League L on S.Id = L.SportId
         join Event E on L.Id = E.LeagueId
         join Line ln on ln.EventId = E.Id
         left join Moneyline M on ln.LineId = M.LineId
