import { env } from "@/env/client";
import { getCentrifugoTokenServerFn } from "@/lib/centrifugo/get-centrifugo-token-server-fn";
import { useServerFn } from "@tanstack/react-start";
import { Centrifuge, PublicationContext, Subscription } from "centrifuge";
import React, {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { Channel } from "./channel";

// Type for the data callback
type DataCallback = (data: any) => void;

interface SubscriptionEntry {
  sub: Subscription;
  count: number;
  listeners: Set<DataCallback>;
  // We keep a reference to the wrapper function for each listener to remove it correctly later
  wrappers: Map<DataCallback, (ctx: PublicationContext) => void>;
}

interface CentrifugoContextType {
  isConnected: boolean;
  client: Centrifuge | null;
  clientId: string | null;
  // Returns a cleanup function to unsubscribe
  subscribe: (channel: Channel, onData: DataCallback) => () => void;
}

const CentrifugoContext = createContext<CentrifugoContextType | null>(null);

interface Props {
  children: ReactNode;
}

export const CentrifugoProvider: React.FC<Props> = ({ children }) => {
  const [client, setClient] = useState<Centrifuge | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [clientId, setClientId] = useState<string | null>(null);
  const getCentrifugoToken = useServerFn(getCentrifugoTokenServerFn);

  // Store active subscriptions and their listener count
  const subsRef = useRef<Record<string, SubscriptionEntry>>({});

  useEffect(() => {
    // Initialize Centrifugo
    const centrifuge = new Centrifuge(env.VITE_WEBSOCKET_URL, {
      getToken: async () => {
        try {
          const res = await getCentrifugoToken();
          return res.data?.token ?? "";
        } catch (error) {
          console.error("Error fetching Centrifugo token:", error);
          return "";
        }
      },
    });

    centrifuge.on("connected", (ctx) => {
      setIsConnected(true);
      setClientId(ctx.client);
      console.log("Centrifugo connected", ctx);
    });

    centrifuge.on("disconnected", (ctx) => {
      setIsConnected(false);
      setClientId(null);
      console.log("Centrifugo disconnected", ctx);
    });

    centrifuge.connect();
    setClient(centrifuge);

    return () => {
      centrifuge.disconnect();
    };
  }, []);

  const subscribe = useCallback(
    (channel: Channel, onData: DataCallback) => {
      if (!client) return () => {};

      // 1. Initialize subscription entry if it doesn't exist
      if (!subsRef.current[channel]) {
        console.log(`[Centrifugo] Creating new subscription: ${channel}`);
        const sub = client.newSubscription(channel);
        sub.subscribe();

        subsRef.current[channel] = {
          sub,
          count: 0,
          listeners: new Set(),
          wrappers: new Map(),
        };
      }

      const entry = subsRef.current[channel];

      // 2. Create a wrapper for the listener (so we can pass just 'ctx.data')
      const listenerWrapper = (ctx: PublicationContext) => {
        onData(ctx.data);
      };

      // 3. Register the listener
      entry.listeners.add(onData);
      entry.wrappers.set(onData, listenerWrapper);
      entry.sub.on("publication", listenerWrapper);
      entry.count += 1;

      // 4. Return Cleanup Function
      return () => {
        const currentEntry = subsRef.current[channel];
        if (!currentEntry) return;

        // Remove the specific listener
        const wrapper = currentEntry.wrappers.get(onData);
        if (wrapper) {
          currentEntry.sub.removeListener("publication", wrapper);
        }
        currentEntry.listeners.delete(onData);
        currentEntry.wrappers.delete(onData);
        currentEntry.count -= 1;

        console.log(
          `[Centrifugo] Unsubscribing component from ${channel}. Remaining: ${currentEntry.count}`
        );

        // If no components are listening anymore, kill the subscription
        if (currentEntry.count <= 0) {
          console.log(`[Centrifugo] Closing real subscription: ${channel}`);
          currentEntry.sub.removeAllListeners();
          currentEntry.sub.unsubscribe();
          client.removeSubscription(currentEntry.sub);
          delete subsRef.current[channel];
        }
      };
    },
    [client]
  );

  return (
    <CentrifugoContext.Provider
      value={{ isConnected, client, clientId, subscribe }}
    >
      {children}
    </CentrifugoContext.Provider>
  );
};

export const useCentrifugo = () => {
  const context = useContext(CentrifugoContext);
  if (!context) {
    throw new Error("useCentrifugo must be used within a CentrifugoProvider");
  }
  return context;
};
