import { useCallback, useEffect, useState } from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { api } from "./api";
import { Layout } from "./components/Layout";
import { InquiryProvider } from "./context/InquiryContext";
import { SiteProvider } from "./context/SiteContext";
import { ArticleDetail } from "./pages/ArticleDetail";
import { Contact } from "./pages/Contact";
import { ContentList } from "./pages/Content";
import { Home } from "./pages/Home";
import { Inquiry } from "./pages/Inquiry";
import { ProductDetail } from "./pages/ProductDetail";
import { Products } from "./pages/Products";
import { SearchPage } from "./pages/Search";
import { StaticPage } from "./pages/StaticPage";
import { VideoPage } from "./pages/Video";
import type { Category, NavItem, Settings } from "./types";

export default function App() {
  const [site, setSite] = useState<{ settings: Settings; categories: Category[]; navigation: NavItem[] }>({ settings: {}, categories: [], navigation: [] });
  const [ready, setReady] = useState(false);
  const [error, setError] = useState("");
  const refreshSite = useCallback(() => api<{ settings: Settings; categories: Category[]; navigation: NavItem[] }>("/bootstrap").then((data) => {
    setSite(data);
    setError("");
  }), []);
  useEffect(() => { refreshSite().catch((e) => setError(e.message)).finally(() => setReady(true)); }, [refreshSite]);
  useEffect(() => {
    function refreshVisibleSite() {
      if (document.visibilityState === "hidden") return;
      refreshSite().catch(() => {});
    }
    window.addEventListener("focus", refreshVisibleSite);
    document.addEventListener("visibilitychange", refreshVisibleSite);
    return () => {
      window.removeEventListener("focus", refreshVisibleSite);
      document.removeEventListener("visibilitychange", refreshVisibleSite);
    };
  }, [refreshSite]);
  useEffect(() => { document.title = site.settings.site_name || ""; }, [site.settings.site_name]);
  if (!ready) return <div className="state-page"><span className="loader" />Loading website...</div>;
  if (error) return <div className="state-page"><h1>Website data is unavailable</h1><p>{error}</p><small>Start the Gin API and initialize the shared database, then refresh this page.</small></div>;
  return <BrowserRouter><SiteProvider value={site}><InquiryProvider><Routes><Route element={<Layout />}><Route index element={<Home />} /><Route path="about-us" element={<StaticPage fixedSlug="about-us" />} /><Route path="page/:slug" element={<StaticPage />} /><Route path="products" element={<Products />} /><Route path="category/:slug" element={<Products />} /><Route path="product/:slug" element={<ProductDetail />} /><Route path="news" element={<ContentList type="news" />} /><Route path="knowledge" element={<ContentList type="knowledge" />} /><Route path="blog" element={<ContentList type="blog" />} /><Route path="article/:slug" element={<ArticleDetail />} /><Route path="contact-us" element={<Contact />} /><Route path="feedback" element={<Inquiry />} /><Route path="video" element={<VideoPage />} /><Route path="search" element={<SearchPage />} /><Route path="*" element={<StaticPage fixedSlug="about-us" />} /></Route></Routes></InquiryProvider></SiteProvider></BrowserRouter>;
}
