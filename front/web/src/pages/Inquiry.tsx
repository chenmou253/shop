import { useState, type FormEvent } from "react";
import { FileUp, ShieldCheck, Trash2 } from "lucide-react";
import { submitForm } from "../api";
import { PageHero } from "../components/Common";
import { useInquiry } from "../context/InquiryContext";

export function Inquiry() {
  const { items, remove, clear } = useInquiry();
  const [message, setMessage] = useState("");
  const [sending, setSending] = useState(false);
  function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault(); setSending(true); setMessage("");
    const form = e.currentTarget; const data = new FormData(form);
    data.set("product_ids", items.map((item) => item.id).join(",")); data.set("source", "front-feedback");
    submitForm("/inquiries", data).then((result) => { setMessage(result.message); clear(); form.reset(); }).catch((error) => setMessage(error.message)).finally(() => setSending(false));
  }
  return <><PageHero title="Send Inquiry" crumbs={["Feedback"]} /><section className="shell inquiry-layout page-space"><div className="inquiry-copy"><span className="eyebrow">Get a tailored solution</span><h1>Tell us what your project needs</h1><p>Upload drawings, sample pictures, a BOM or specifications. Our sales and engineering team will respond within 24 hours.</p><div className="trust-note"><ShieldCheck /><span><b>Your information stays confidential.</b> We use it only to prepare your quotation and technical recommendation.</span></div>{items.length > 0 && <div className="inquiry-products"><h3>Selected Products ({items.length}/10)</h3>{items.map((item) => <div key={item.id}><img src={item.main_image} alt="" /><span>{item.name}</span><button title="Remove product" onClick={() => remove(item.id)}><Trash2 size={17} /></button></div>)}</div>}</div><form className="inquiry-form" onSubmit={submit}><label><span>Your Name *</span><input name="name" required /></label><label><span>E-mail *</span><input name="email" type="email" required /></label><label><span>Phone / WhatsApp</span><input name="phone" /></label><label className="wide"><span>Message *</span><textarea name="message" rows={7} required placeholder="Please include material, size, quantity, standards and delivery country if known." /></label><label className="upload-field wide"><FileUp /><span>Send a file or picture<small>PDF, Office, JPG, PNG or WebP, up to 10 MB</small></span><input name="attachment" type="file" accept=".pdf,.doc,.docx,.xls,.xlsx,.jpg,.jpeg,.png,.webp" /></label><button className="button large wide" disabled={sending}>{sending ? "Sending..." : "Submit Inquiry"}</button>{message && <p className="form-message wide">{message}</p>}</form></section></>;
}
