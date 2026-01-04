"use client";

import { useEffect } from "react";
import { useAuth } from "@clerk/nextjs";
import { setTokenGetter } from "@/lib/api";

/**
 * Hook to connect Clerk authentication with the API client.
 * This should be called once at the app level to set up token injection.
 */
export function useAuthenticatedApi() {
  const { getToken, isLoaded, isSignedIn } = useAuth();

  useEffect(() => {
    if (isLoaded) {
      // Set up the token getter for the API client
      setTokenGetter(async () => {
        if (!isSignedIn) {
          return null;
        }
        return getToken();
      });
    }
  }, [getToken, isLoaded, isSignedIn]);

  return { isLoaded, isSignedIn };
}
