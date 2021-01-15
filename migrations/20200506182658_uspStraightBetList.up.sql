create or alter proc dbo.uspStraightBetList @TVP dbo.StraightBetListType READONLY as
begin
    set nocount on

    MERGE dbo.StraightBetList AS t
    USING @TVP s
    ON (t.BetId = s.BetId)

    WHEN MATCHED THEN
        UPDATE
        SET BetId              = s.BetId,
            WagerNumber        = s.WagerNumber,
            PlacedAt           = s.PlacedAt,
            BetStatus          = s.BetStatus,
            BetType            = s.BetType,
            Win                = s.Win,
            Risk               = s.Risk,
            WinLoss            = s.WinLoss,
            OddsFormat         = s.OddsFormat,
            CustomerCommission = s.CustomerCommission,
--             CancellationReason = s.CancellationReason,
            UpdateSequence     = s.UpdateSequence,
            SportId            = s.SportId,
            LeagueId           = s.LeagueId,
            EventId            = s.EventId,
            Handicap           = s.Handicap,
            Price              = s.Price,
            TeamName           = s.TeamName,
            Side               = s.Side,
            Pitcher1           = s.Pitcher1,
            Pitcher2           = s.Pitcher2,
            Pitcher1MustStart  = s.Pitcher1MustStart,
            Pitcher2MustStart  = s.Pitcher2MustStart,
            Team1              = s.Team1,
            Team2              = s.Team2,
            PeriodNumber       = s.PeriodNumber,
            Team1Score         = s.Team1Score,
            Team2Score         = s.Team2Score,
            FtTeam1Score       = s.FtTeam1Score,
            FtTeam2Score       = s.FtTeam2Score,
            PTeam1Score        = s.PTeam1Score,
            PTeam2Score        = s.PTeam2Score,
            EventStartTime     = s.EventStartTime,
            UpdatedAt          =sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (BetId, WagerNumber, PlacedAt, BetStatus, BetType, Win, Risk, WinLoss, OddsFormat, CustomerCommission,
--                 CancellationReason,
                UpdateSequence, SportId, LeagueId, EventId, Handicap, Price, TeamName, Side, Pitcher1, Pitcher2,
                Pitcher1MustStart, Pitcher2MustStart, Team1, Team2, PeriodNumber, Team1Score, Team2Score, FtTeam1Score,
                FtTeam2Score, PTeam1Score, PTeam2Score, EventStartTime)
        VALUES (s.BetId, s.WagerNumber, s.PlacedAt, s.BetStatus, s.BetType, s.Win, s.Risk, s.WinLoss, s.OddsFormat,
                s.CustomerCommission,
--                 s.CancellationReason,
                s.UpdateSequence, s.SportId, s.LeagueId, s.EventId, s.Handicap, s.Price, s.TeamName, s.Side,
                s.Pitcher1, s.Pitcher2,
                s.Pitcher1MustStart, s.Pitcher2MustStart, s.Team1, s.Team2, s.PeriodNumber, s.Team1Score, s.Team2Score,
                s.FtTeam1Score,
                s.FtTeam2Score, s.PTeam1Score, s.PTeam2Score, s.EventStartTime);
end