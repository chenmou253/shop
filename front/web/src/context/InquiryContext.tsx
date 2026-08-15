import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from "react";
import type { Product } from "../types";

type InquiryValue = {
  items: Product[];
  add: (product: Product) => void;
  remove: (id: number) => void;
  clear: () => void;
};

const InquiryContext = createContext<InquiryValue | null>(null);

export function InquiryProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<Product[]>(() => {
    try { return JSON.parse(localStorage.getItem("depo-inquiry") || "[]"); } catch { return []; }
  });

  useEffect(() => localStorage.setItem("depo-inquiry", JSON.stringify(items)), [items]);

  const value = useMemo(() => ({
    items,
    add(product: Product) {
      setItems((current) => current.some((item) => item.id === product.id) || current.length >= 10 ? current : [...current, product]);
    },
    remove(id: number) { setItems((current) => current.filter((item) => item.id !== id)); },
    clear() { setItems([]); }
  }), [items]);

  return <InquiryContext.Provider value={value}>{children}</InquiryContext.Provider>;
}

export function useInquiry() {
  const value = useContext(InquiryContext);
  if (!value) throw new Error("useInquiry must be used inside InquiryProvider");
  return value;
}
