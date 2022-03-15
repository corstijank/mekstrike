FROM mcr.microsoft.com/dotnet/sdk:6.0
WORKDIR /app

# Copy csproj and restore as distinct layers
COPY BattlefieldService BattlefieldService
COPY Battlefield.Interfaces Battlefield.Interfaces

WORKDIR /app/BattlefieldService
RUN dotnet restore

# Copy everything else and build
RUN dotnet publish -c Release 

CMD ["/app/BattlefieldService/bin/Release/net6.0/BattlefieldService"]