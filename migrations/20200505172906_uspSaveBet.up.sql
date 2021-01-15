create or alter proc dbo.uspSaveBet @SurebetId bigint,
                                    @SideIndex tinyint,
                                    @BetId bigint = null,
                                    @TryCount tinyint= null,
                                    @Status varchar(1000)= null,
                                    @StatusInfo varchar(1000)= null,
                                    @Start bigint= null,
                                    @Done bigint= null,
                                    @Price decimal(9, 5)= null,
                                    @Stake decimal(9, 5)= null,
                                    @ApiBetId bigint= null
as
begin
    set nocount on

    MERGE dbo.Bet AS t
    USING (select @SurebetId,
                  @SideIndex,
                  @BetId,
                  @TryCount,
                  @Status,
                  @StatusInfo,
                  @Start,
                  @Done,
                  @Price,
                  @Stake,
                  @ApiBetId) s (SurebetId, SideIndex, BetId, TryCount, Status, StatusInfo, Start, Done, Price, Stake, ApiBetId)
    ON (t.SurebetId = s.SurebetId and t.SideIndex=s.SideIndex)

    WHEN MATCHED THEN
        UPDATE
        SET BetId      = s.BetId,
            TryCount   = s.TryCount,
            Status     = s.Status,
            StatusInfo = s.StatusInfo,
            Start      = s.Start,
            Done       = s.Done,
            Price      = s.Price,
            Stake      = s.Stake,
            ApiBetId   = s.ApiBetId,
            UpdatedAt  =sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (SurebetId, SideIndex,BetId, TryCount, Status, StatusInfo, Start, Done, Price, Stake, ApiBetId)
        VALUES (s.SurebetId, s.SideIndex, s.BetId, s.TryCount, s.Status, s.StatusInfo, s.Start, s.Done, s.Price, s.Stake, s.ApiBetId);
end