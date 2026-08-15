import { useEffect, useState } from "react";
import { ArrowLeft, ArrowRight, Check, Mail } from "lucide-react";
import { Link } from "react-router-dom";
import { api, externalLink } from "../api";
import { ProductCard, SectionTitle } from "../components/Common";
import { useSite } from "../context/SiteContext";
import type { HomeData } from "../types";

export function Home() {
  const { settings } = useSite();
  const [data, setData] = useState<HomeData | null>(null);
  const [error, setError] = useState("");
  useEffect(() => { api<HomeData>("/home").then(setData).catch((e) => setError(e.message)); }, []);
  if (error) return <div className="state-page">{error}</div>;
  if (!data) return <div className="state-page"><span className="loader" />Loading...</div>;
  const homeSettings = { ...data.settings, ...settings };
  const whyStyle = homeSettings.why_background_url ? { backgroundImage: `linear-gradient(rgba(0, 76, 136, .82), rgba(0, 76, 136, .82)), url(${homeSettings.why_background_url})` } : undefined;

  return <>
    <HeroSlider banners={data.banners} />
    <section className="featured-categories shell">
      {data.featured_categories.map((category) => <article className="category-feature" key={category.id}>
        <Link to={`/category/${category.slug}`}><img src={category.image_url} alt={category.name} /></Link>
        <div><h2>{category.name}</h2><p>{category.summary}</p><ul className="category-mobile-list">{category.summary.split(/[,;]/).map((item) => item.trim()).filter(Boolean).slice(0, 3).map((item, index) => <li key={`${item}-${index}`}>{item}</li>)}</ul><Link to={`/category/${category.slug}`}>View Products <ArrowRight size={17} /></Link></div>
      </article>)}
    </section>
    <section className="section hot-section"><div className="shell"><SectionTitle>Hot Products</SectionTitle><div className="product-grid home-products">{data.hot_products.map((product) => <ProductCard key={product.id} product={product} />)}</div><div className="center"><Link className="button pill" to="/products">Learn More <ArrowRight size={18} /></Link></div></div></section>
    <section className="why-section" style={whyStyle}><div className="shell"><SectionTitle>Why Depo</SectionTitle><div className="why-grid"><div className="benefit-list">{data.benefits.map((benefit) => <p key={benefit.id}><Check size={18} />{benefit.title}</p>)}</div><VideoFrame url={homeSettings.home_video_url} /></div></div></section>
    <section className="section partners-section"><div className="shell"><SectionTitle>Partners</SectionTitle><div className="partner-grid">{data.partners.map((partner) => <a href={partner.website_url || undefined} key={partner.id} title={partner.name}><img src={partner.logo_url} alt={partner.name} loading="lazy" /></a>)}</div></div></section>
    <section className="section team-section"><div className="shell"><SectionTitle>Meet Our Team</SectionTitle><div className="team-grid">{data.team.map((member) => <article className="team-card" key={member.id}><img src={member.image_url} alt={member.name} /><div className="team-caption"><h3>{member.name}</h3><p>{member.role}</p></div><a href={`mailto:${member.email}`}><Mail size={18} />{member.email}</a></article>)}</div><div className="center"><Link className="button pill" to="/about-us">View More</Link></div></div></section>
    <Industries items={data.industries} />
    <section className="section news-section"><div className="shell"><SectionTitle>News</SectionTitle><div className="news-grid">{data.news.map((article) => <article className="news-card" key={article.id}><Link to={`/article/${article.slug}`}><img src={article.cover_image} alt={article.title} /></Link><div><time>{new Date(article.published_at).toLocaleDateString("en-US", { month: "short", day: "2-digit", year: "numeric" })}</time><h3><Link to={`/article/${article.slug}`}>{article.title}</Link></h3><p>{article.summary}</p><Link className="more-link" to={`/article/${article.slug}`}>Learn More <ArrowRight size={16} /></Link></div></article>)}</div><div className="center"><Link className="button pill" to="/news">View More</Link></div></div></section>
  </>;
}

function HeroSlider({ banners }: { banners: HomeData["banners"] }) {
  const [active, setActive] = useState(0);
  useEffect(() => {
    if (banners.length < 2) return;
    const timer = window.setInterval(() => setActive((n) => (n + 1) % banners.length), 5500);
    return () => clearInterval(timer);
  }, [banners.length]);
  if (!banners.length) return null;
  const move = (n: number) => setActive((active + n + banners.length) % banners.length);
  return <section className="hero-slider">{banners.map((banner, index) => <Link className={`hero-slide ${index === active ? "active" : ""}`} key={banner.id} to={externalLink(banner.link_url || "/products")}><img src={banner.image_url} alt={banner.title} /><div className="hero-copy shell"><h1>{banner.title}</h1><p>{banner.subtitle}</p></div></Link>)}<button className="hero-prev" title="Previous slide" onClick={() => move(-1)}><ArrowLeft /></button><button className="hero-next" title="Next slide" onClick={() => move(1)}><ArrowRight /></button><div className="hero-dots">{banners.map((b, index) => <button className={active === index ? "active" : ""} onClick={() => setActive(index)} key={b.id} title={`Slide ${index + 1}`} />)}</div></section>;
}

function VideoFrame({ url }: { url?: string }) {
  const id = url?.match(/(?:youtu\.be\/|youtube\.com\/(?:watch\?v=|embed\/))([^?&/]+)/)?.[1];
  if (!id) return null;
  return <div className="video-frame"><iframe src={`https://www.youtube-nocookie.com/embed/${id}`} title="Depo factory video" allowFullScreen loading="lazy" /></div>;
}

function Industries({ items }: { items: HomeData["industries"] }) {
  const [active, setActive] = useState(0);
  if (!items.length) return null;
  const item = items[Math.min(active, items.length - 1)];
  return <section className="section industries-section"><div className="shell"><SectionTitle>Products Industry Application</SectionTitle><div className="industry-tabs">{items.map((industry, index) => <button className={index === active ? "active" : ""} key={industry.id} onClick={() => setActive(index)}>{industry.title}</button>)}</div><article className="industry-panel"><div><h3>{item.title}</h3><h4>{item.subtitle}</h4><p>{item.description}</p><Link className="button pill" to={externalLink(item.link_url || "/products")}>View More <ArrowRight size={18} /></Link></div><img src={item.image_url} alt={item.title} /></article></div></section>;
}
