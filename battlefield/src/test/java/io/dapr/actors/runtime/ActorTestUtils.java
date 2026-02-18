package io.dapr.actors.runtime;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import io.dapr.actors.ActorId;
import reactor.core.publisher.Mono;

/**
 * Test utility for creating mocked Dapr actor contexts. Lives in io.dapr.actors.runtime to access
 * package-private types.
 */
public class ActorTestUtils {

  @SuppressWarnings("unchecked")
  public static <T extends AbstractActor> ActorRuntimeContext<T> createMockContext(
      String actorTypeName) {
    var stateProvider = mock(DaprStateAsyncProvider.class);
    when(stateProvider.contains(anyString(), any(ActorId.class), anyString()))
        .thenReturn(Mono.just(false));

    var typeInfo = mock(ActorTypeInformation.class);
    when(typeInfo.getName()).thenReturn(actorTypeName);

    var context = mock(ActorRuntimeContext.class);
    when(context.getStateProvider()).thenReturn(stateProvider);
    when(context.getActorTypeInformation()).thenReturn(typeInfo);

    return context;
  }
}
