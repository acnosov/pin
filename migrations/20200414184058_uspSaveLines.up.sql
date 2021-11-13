create or alter proc dbo.uspSaveLines @TVP dbo.LineType READONLY as
begin
    set nocount on

    MERGE dbo.Line AS t
    USING @TVP s
    ON (t.LineId = s.LineId)

    WHEN MATCHED THEN
        UPDATE
        SET EventId      = s.EventId,
            Number       = s.Number,
            Cutoff       = s.Cutoff,
            Status       = s.Status,
            MaxSpread    = s.MaxSpread,
            MaxMoneyline = s.MaxMoneyline,
            MaxTotal     = s.MaxTotal,
            MaxTeamTotal = s.MaxTeamTotal,
            UpdatedAt    = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (LineId,
                EventId,
                Number,
                Cutoff,
                Status,
                MaxSpread,
                MaxMoneyline,
                MaxTotal,
                MaxTeamTotal)
        VALUES (s.LineId,
                s.EventId,
                s.Number,
                s.Cutoff,
                s.Status,
                s.MaxSpread,
                s.MaxMoneyline,
                s.MaxTotal,
                s.MaxTeamTotal);
end
