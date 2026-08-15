import { createContext, useContext, type ReactNode } from "react";
import type { Category, NavItem, Settings } from "../types";

type SiteValue = { settings: Settings; categories: Category[]; navigation: NavItem[] };
const SiteContext = createContext<SiteValue>({ settings: {}, categories: [], navigation: [] });

export function SiteProvider({ value, children }: { value: SiteValue; children: ReactNode }) {
  return <SiteContext.Provider value={value}>{children}</SiteContext.Provider>;
}

export function useSite() { return useContext(SiteContext); }
