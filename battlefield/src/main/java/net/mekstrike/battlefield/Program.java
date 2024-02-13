package net.mekstrike.battlefield;

import java.io.IOException;
import java.time.Duration;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import io.dapr.actors.runtime.ActorRuntime;
import io.dapr.serializer.DefaultObjectSerializer;
import net.mekstrike.serialization.MekstrikeSerializer;

@SpringBootApplication
public class Program {
    public static void main(String[] args) throws IOException, InterruptedException {
       new Program().start();
    }

    public void start() throws IOException,InterruptedException{
         // Idle timeout until actor instance is deactivated.
         ActorRuntime.getInstance().getConfig().setActorIdleTimeout(Duration.ofSeconds(30));
         // How often actor instances are scanned for deactivation and balance.
         ActorRuntime.getInstance().getConfig().setActorScanInterval(Duration.ofSeconds(10));
         // How long to wait until for draining an ongoing API call for an actor
         // instance.
         ActorRuntime.getInstance().getConfig().setDrainOngoingCallTimeout(Duration.ofSeconds(10));
         // Determines whether to drain API calls for actors instances being balanced.
         ActorRuntime.getInstance().getConfig().setDrainBalancedActors(true);
         // Register the Actor class.
         ActorRuntime.getInstance().registerActor(BattlefieldImpl.class, new MekstrikeSerializer(),new DefaultObjectSerializer());
         // Start Dapr's callback endpoint.
         SpringApplication app = new SpringApplication(Program.class);
         app.run(String.format("--server.port=%d", 9000));
    }
}
