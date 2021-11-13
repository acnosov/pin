create or alter function dbo.fnFindSportIdByName(@SportName varchar(180)) RETURNS int as
begin
    declare @SportId int
    select @SportId = Id from dbo.Sport where Name = @SportName
    return @SportId
end

create or alter function dbo.fnFindLeagueIdByName(@Name varchar(180), @SportId int) RETURNS int as
begin
    declare @Id int
    select @Id = Id from dbo.League where  Name = @Name and SportId = @SportId
    return @Id
end


create or alter function dbo.fnFindEvent(@Starts datetimeoffset, @Home varchar(300), @Away varchar(300), @LeagueId int) RETURNS int as
begin
    declare @Id int
    select @Id = Id from dbo.Event where Starts = @Starts and Home = @Home and Away = @Away and LeagueId = @LeagueId
    return @Id
end


