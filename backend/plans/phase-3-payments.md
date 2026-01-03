# Phase 3: Payment & Upsells

**Sprint:** 3 (Weeks 5-6)  
**Focus:** Checkout integration and one-click upsell flow  
**Status:** Ready for Implementation

---

## Overview

Phase 3 implements the core monetization features: inline checkout with PIX payment and one-click upsell functionality. This phase validates the central hypothesis that fluid upsells increase AOV without adding friction.

---

## Objectives

1. Integrate Mercado Pago for PIX payments
2. Build inline checkout component for sales pages
3. Implement one-click upsell flow
4. Handle payment webhooks and confirmations
5. Ensure secure payment data handling

---

## Features

### Feature 3: Checkout Inline com PIX

**Priority:** CRITICAL  
**Complexity:** High  
**Estimated Time:** 2 weeks

#### Description

Embedded checkout form directly in the sales page that accepts PIX payments through Mercado Pago, eliminating the need for external checkout pages and reducing friction.

#### Key Requirements

**Checkout Form:**
- Name, email, phone fields with validation
- PIX payment method (MVP only)
- Real-time form validation
- Mobile-optimized input types

**PIX Payment Flow:**
1. User fills form and clicks "Finalizar Compra"
2. Backend creates payment in Mercado Pago
3. QR Code displayed to user
4. Frontend polls for payment status
5. On confirmation, redirect to upsell

**Security:**
- PCI DSS compliance (handled by Mercado Pago)
- Secure webhook validation
- Payment data encryption in transit

#### Implementation Tasks

**1. Mercado Pago Integration**

**Task 1.1: Mercado Pago Configuration** (`/backend/src/config/mercadopago.config.ts`)
```typescript
import mercadopago from 'mercadopago';

mercadopago.configure({
  access_token: process.env.MERCADOPAGO_ACCESS_TOKEN!
});

export default mercadopago;
```

**Task 1.2: Payment Service** (`/backend/src/services/payment.service.ts`)
```typescript
import mercadopago from '../config/mercadopago.config';

export interface CreatePaymentDTO {
  amount: number;
  description: string;
  payer: {
    name: string;
    email: string;
    phone: string;
  };
  metadata: {
    funnelId: string;
    sessionId: string;
  };
}

export interface PaymentResult {
  paymentId: string;
  qrCode: string;
  qrCodeBase64: string;
  expiresAt: Date;
}

export class PaymentService {
  async createPixPayment(data: CreatePaymentDTO): Promise<PaymentResult> {
    const payment = await mercadopago.payment.create({
      transaction_amount: data.amount,
      description: data.description,
      payment_method_id: 'pix',
      payer: {
        email: data.payer.email,
        first_name: data.payer.name.split(' ')[0],
        last_name: data.payer.name.split(' ').slice(1).join(' ') || data.payer.name.split(' ')[0]
      },
      notification_url: `${process.env.API_URL}/webhooks/mercadopago`,
      metadata: data.metadata
    });

    if (payment.status !== 'pending') {
      throw new Error('Payment creation failed');
    }

    const qrData = payment.body.point_of_interaction.transaction_data;

    return {
      paymentId: payment.body.id.toString(),
      qrCode: qrData.qr_code,
      qrCodeBase64: qrData.qr_code_base64,
      expiresAt: new Date(Date.now() + 15 * 60 * 1000) // 15 minutes
    };
  }

  async getPaymentStatus(paymentId: string): Promise<string> {
    const payment = await mercadopago.payment.get(paymentId);
    return payment.body.status;
  }

  async processPaymentApproval(paymentId: string): Promise<void> {
    const payment = await mercadopago.payment.get(paymentId);
    
    if (payment.body.status !== 'approved') {
      throw new Error('Payment not approved');
    }

    const metadata = payment.body.metadata;
    const sessionId = metadata.session_id;

    // Update session state
    await sessionService.updateSessionState(sessionId, 'purchased');

    // Save payment info to session for one-click upsell
    await sessionService.savePaymentData(sessionId, {
      paymentId: payment.body.id,
      payerId: payment.body.payer.id,
      paymentMethod: 'pix',
      amount: payment.body.transaction_amount,
      customerEmail: payment.body.payer.email
    });

    // Create order record
    await orderService.create({
      funnelId: metadata.funnel_id,
      sessionId: sessionId,
      paymentId: payment.body.id.toString(),
      amount: payment.body.transaction_amount,
      status: 'paid',
      customerName: `${payment.body.payer.first_name} ${payment.body.payer.last_name}`,
      customerEmail: payment.body.payer.email
    });
  }
}

export const paymentService = new PaymentService();
```

**Task 1.3: Order Service** (`/backend/src/services/order.service.ts`)
```typescript
export interface Order {
  id: string;
  funnelId: string;
  sessionId: string;
  paymentId: string;
  amount: number;
  status: 'pending' | 'paid' | 'refunded';
  customerName: string;
  customerEmail: string;
  items: OrderItem[];
  createdAt: Date;
  updatedAt: Date;
}

export interface OrderItem {
  productId: string;
  productName: string;
  amount: number;
  type: 'main' | 'upsell';
}

export class OrderService {
  async create(data: Partial<Order>): Promise<Order> {
    return await db.orders.create({
      data: {
        ...data,
        items: [
          {
            productId: data.funnelId!,
            productName: 'Main Product',
            amount: data.amount!,
            type: 'main'
          }
        ]
      }
    });
  }

  async addUpsellItem(orderId: string, item: OrderItem): Promise<Order> {
    const order = await db.orders.findUnique({ where: { id: orderId } });
    
    return await db.orders.update({
      where: { id: orderId },
      data: {
        items: [...order.items, item],
        amount: order.amount + item.amount
      }
    });
  }

  async getBySessionId(sessionId: string): Promise<Order | null> {
    return await db.orders.findFirst({
      where: { sessionId }
    });
  }
}

export const orderService = new OrderService();
```

**2. Checkout Component**

**Task 2.1: Checkout Form Component** (`/frontend/src/components/CheckoutInline.tsx`)
```typescript
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';

const checkoutSchema = z.object({
  name: z.string().min(3, 'Nome deve ter pelo menos 3 caracteres'),
  email: z.string().email('Email inválido'),
  phone: z.string().regex(/^\(\d{2}\) \d{5}-\d{4}$/, 'Telefone inválido')
});

type CheckoutFormData = z.infer<typeof checkoutSchema>;

export function CheckoutInline({ amount, productName }: CheckoutInlineProps) {
  const [isProcessing, setIsProcessing] = useState(false);
  const [paymentData, setPaymentData] = useState<PaymentResult | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<CheckoutFormData>({
    resolver: zodResolver(checkoutSchema)
  });

  const onSubmit = async (data: CheckoutFormData) => {
    setIsProcessing(true);

    try {
      const response = await fetch('/api/payments/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          amount,
          description: productName,
          payer: data
        })
      });

      const result = await response.json();
      setPaymentData(result);
      
      // Start polling for payment status
      startPaymentPolling(result.paymentId);
    } catch (error) {
      console.error('Payment creation failed:', error);
      alert('Erro ao processar pagamento. Tente novamente.');
    } finally {
      setIsProcessing(false);
    }
  };

  if (paymentData) {
    return (
      <PixPaymentDisplay
        paymentData={paymentData}
        onPaymentConfirmed={() => {
          window.location.href = '/upsell';
        }}
      />
    );
  }

  return (
    <div className="checkout-inline max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg">
      <h2 className="text-2xl font-bold mb-6">Checkout</h2>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">
            Nome Completo
          </label>
          <input
            {...register('name')}
            type="text"
            placeholder="Seu nome completo"
            className="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500"
          />
          {errors.name && (
            <p className="text-red-500 text-sm mt-1">{errors.name.message}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">
            Email
          </label>
          <input
            {...register('email')}
            type="email"
            placeholder="seu@email.com"
            className="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500"
          />
          {errors.email && (
            <p className="text-red-500 text-sm mt-1">{errors.email.message}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">
            Telefone
          </label>
          <input
            {...register('phone')}
            type="tel"
            placeholder="(11) 99999-9999"
            className="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500"
          />
          {errors.phone && (
            <p className="text-red-500 text-sm mt-1">{errors.phone.message}</p>
          )}
        </div>

        <div className="bg-gray-50 p-4 rounded-lg">
          <div className="flex items-center space-x-2">
            <input
              type="radio"
              id="pix"
              checked
              readOnly
              className="w-4 h-4"
            />
            <label htmlFor="pix" className="flex items-center space-x-2">
              <span className="font-medium">PIX</span>
              <span className="text-sm text-gray-600">(Pagamento Instantâneo)</span>
            </label>
          </div>
        </div>

        <div className="text-center py-4">
          <p className="text-3xl font-bold">
            R$ {amount.toFixed(2).replace('.', ',')}
          </p>
        </div>

        <button
          type="submit"
          disabled={isProcessing}
          className="w-full bg-green-600 text-white py-4 rounded-lg font-bold text-lg hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {isProcessing ? 'Processando...' : 'FINALIZAR COMPRA'}
        </button>

        <div className="text-center text-sm text-gray-600">
          🔒 Compra 100% Segura
        </div>
      </form>
    </div>
  );
}
```

**Task 2.2: PIX Payment Display** (`/frontend/src/components/PixPaymentDisplay.tsx`)
```typescript
export function PixPaymentDisplay({
  paymentData,
  onPaymentConfirmed
}: PixPaymentDisplayProps) {
  const [timeLeft, setTimeLeft] = useState(15 * 60); // 15 minutes in seconds
  const [status, setStatus] = useState<'pending' | 'approved' | 'expired'>('pending');

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 0) {
          clearInterval(timer);
          setStatus('expired');
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  const copyToClipboard = async () => {
    await navigator.clipboard.writeText(paymentData.qrCode);
    toast.success('Código PIX copiado!');
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  if (status === 'approved') {
    return (
      <div className="text-center p-8">
        <div className="text-6xl mb-4">✅</div>
        <h2 className="text-2xl font-bold mb-2">Pagamento Confirmado!</h2>
        <p className="text-gray-600 mb-4">Redirecionando...</p>
      </div>
    );
  }

  if (status === 'expired') {
    return (
      <div className="text-center p-8">
        <div className="text-6xl mb-4">⏱️</div>
        <h2 className="text-2xl font-bold mb-2">Tempo Expirado</h2>
        <p className="text-gray-600 mb-4">O código PIX expirou. Tente novamente.</p>
        <button
          onClick={() => window.location.reload()}
          className="bg-blue-600 text-white px-6 py-3 rounded-lg"
        >
          Tentar Novamente
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg">
      <h2 className="text-2xl font-bold mb-4 text-center">Pagamento via PIX</h2>

      <div className="bg-white p-4 rounded-lg mb-4 flex justify-center">
        <img
          src={`data:image/png;base64,${paymentData.qrCodeBase64}`}
          alt="QR Code PIX"
          className="w-64 h-64"
        />
      </div>

      <p className="text-center text-gray-700 mb-4">
        Escaneie o QR Code com o app do seu banco
      </p>

      <div className="text-center text-gray-500 mb-4">OU</div>

      <div className="mb-4">
        <label className="block text-sm font-medium mb-2">
          PIX Copia e Cola:
        </label>
        <div className="flex space-x-2">
          <input
            type="text"
            value={paymentData.qrCode}
            readOnly
            className="flex-1 px-3 py-2 border rounded-lg bg-gray-50 text-sm"
          />
          <button
            onClick={copyToClipboard}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
          >
            Copiar
          </button>
        </div>
      </div>

      <div className="text-center py-4 bg-yellow-50 rounded-lg mb-4">
        <p className="text-sm text-gray-600 mb-1">⏱️ Aguardando pagamento...</p>
        <p className="text-lg font-bold">Válido por: {formatTime(timeLeft)}</p>
      </div>

      <div className="text-center">
        <p className="text-2xl font-bold">
          R$ {paymentData.amount?.toFixed(2).replace('.', ',')}
        </p>
      </div>
    </div>
  );
}
```

**Task 2.3: Payment Polling Hook** (`/frontend/src/hooks/usePaymentPolling.ts`)
```typescript
export function usePaymentPolling(
  paymentId: string,
  onApproved: () => void,
  onFailed: () => void
) {
  useEffect(() => {
    const pollInterval = setInterval(async () => {
      try {
        const response = await fetch(`/api/payments/${paymentId}/status`);
        const data = await response.json();

        if (data.status === 'approved') {
          clearInterval(pollInterval);
          onApproved();
        } else if (data.status === 'rejected' || data.status === 'cancelled') {
          clearInterval(pollInterval);
          onFailed();
        }
      } catch (error) {
        console.error('Payment polling error:', error);
      }
    }, 3000); // Poll every 3 seconds

    return () => clearInterval(pollInterval);
  }, [paymentId, onApproved, onFailed]);
}
```

**3. Backend Routes and Controllers**

**Task 3.1: Payment Routes** (`/backend/src/routes/payment.routes.ts`)
```typescript
import { Router } from 'express';
import { paymentController } from '../controllers/payment.controller';

const router = Router();

router.post('/payments/create', paymentController.createPayment);
router.get('/payments/:paymentId/status', paymentController.getStatus);

export default router;
```

**Task 3.2: Payment Controller** (`/backend/src/controllers/payment.controller.ts`)
```typescript
export class PaymentController {
  async createPayment(req: Request, res: Response) {
    try {
      const { amount, description, payer } = req.body;
      const sessionId = req.cookies.funnel_session_id;
      const session = await sessionService.getSession(sessionId);

      if (!session) {
        return res.status(401).json({ error: 'Invalid session' });
      }

      const result = await paymentService.createPixPayment({
        amount,
        description,
        payer,
        metadata: {
          funnelId: session.funnelId,
          sessionId: session.id
        }
      });

      res.json(result);
    } catch (error) {
      console.error('Payment creation error:', error);
      res.status(500).json({ error: 'Payment creation failed' });
    }
  }

  async getStatus(req: Request, res: Response) {
    try {
      const { paymentId } = req.params;
      const status = await paymentService.getPaymentStatus(paymentId);

      res.json({ status });
    } catch (error) {
      console.error('Payment status error:', error);
      res.status(500).json({ error: 'Failed to get payment status' });
    }
  }
}

export const paymentController = new PaymentController();
```

**4. Webhook Handler**

**Task 4.1: Mercado Pago Webhook** (`/backend/src/routes/webhook.routes.ts`)
```typescript
import { Router } from 'express';
import { webhookController } from '../controllers/webhook.controller';

const router = Router();

router.post('/webhooks/mercadopago', webhookController.handleMercadoPago);

export default router;
```

**Task 4.2: Webhook Controller** (`/backend/src/controllers/webhook.controller.ts`)
```typescript
export class WebhookController {
  async handleMercadoPago(req: Request, res: Response) {
    try {
      const { type, data } = req.body;

      // Acknowledge receipt immediately
      res.status(200).send('OK');

      // Process webhook asynchronously
      if (type === 'payment') {
        const paymentId = data.id;
        
        // Get payment details
        const payment = await mercadopago.payment.get(paymentId);
        
        if (payment.body.status === 'approved') {
          await paymentService.processPaymentApproval(paymentId.toString());
        }
      }
    } catch (error) {
      console.error('Webhook processing error:', error);
      // Still return 200 to prevent retries
      res.status(200).send('OK');
    }
  }
}

export const webhookController = new WebhookController();
```

---

### Feature 4: One-Click Upsell

**Priority:** CRITICAL  
**Complexity:** Medium  
**Estimated Time:** 1 week

#### Description

Dedicated upsell page that appears after main purchase, allowing customers to add complementary products with a single click, without re-entering payment information.

#### Key Requirements

**Upsell Page:**
- Congratulations message
- Product image and benefits
- Price with discount
- Accept/Decline buttons
- Countdown timer (10 minutes)

**One-Click Flow:**
- Use saved payment data from session
- Generate new PIX QR code for upsell amount
- No form required
- Instant processing

**Note:** Since PIX doesn't support tokenization, "one-click" means no form re-entry, but user still needs to scan new QR code.

#### Implementation Tasks

**Task 1: Upsell Page Component** (`/frontend/src/pages/UpsellPage.tsx`)
```typescript
export function UpsellPage() {
  const [timeLeft, setTimeLeft] = useState(10 * 60); // 10 minutes
  const [isProcessing, setIsProcessing] = useState(false);
  const [paymentData, setPaymentData] = useState<PaymentResult | null>(null);

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 0) {
          clearInterval(timer);
          window.location.href = '/confirmation';
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  const handleAccept = async () => {
    setIsProcessing(true);

    try {
      const response = await fetch('/api/upsell/accept', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }
      });

      const result = await response.json();
      setPaymentData(result);
      
      // Start polling for payment
      startPaymentPolling(result.paymentId);
    } catch (error) {
      console.error('Upsell acceptance failed:', error);
      alert('Erro ao processar upsell. Tente novamente.');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleDecline = async () => {
    await fetch('/api/upsell/decline', { method: 'POST' });
    window.location.href = '/confirmation';
  };

  if (paymentData) {
    return (
      <PixPaymentDisplay
        paymentData={paymentData}
        onPaymentConfirmed={() => {
          window.location.href = '/confirmation';
        }}
      />
    );
  }

  return (
    <div className="max-w-2xl mx-auto p-6">
      <div className="text-center mb-8">
        <div className="text-6xl mb-4">🎉</div>
        <h1 className="text-3xl font-bold mb-2">PARABÉNS!</h1>
        <p className="text-xl text-gray-600">Sua compra foi confirmada!</p>
      </div>

      <div className="border-t border-b py-8 mb-8">
        <h2 className="text-2xl font-bold text-center mb-6">
          OFERTA ESPECIAL SÓ PARA VOCÊ
        </h2>

        <div className="mb-6">
          <img
            src="/upsell-product.jpg"
            alt="Produto Upsell"
            className="w-full rounded-lg"
          />
        </div>

        <h3 className="text-xl font-bold mb-4">
          E-book de Receitas Saudáveis
        </h3>

        <ul className="space-y-2 mb-6">
          <li className="flex items-start">
            <span className="text-green-600 mr-2">✓</span>
            50 receitas exclusivas
          </li>
          <li className="flex items-start">
            <span className="text-green-600 mr-2">✓</span>
            Plano alimentar 30 dias
          </li>
          <li className="flex items-start">
            <span className="text-green-600 mr-2">✓</span>
            Lista de compras pronta
          </li>
        </ul>

        <div className="text-center mb-6">
          <p className="text-gray-500 line-through">De R$ 97,00</p>
          <p className="text-4xl font-bold text-green-600">R$ 47,00</p>
          <p className="text-sm text-gray-600">(51% OFF - só hoje!)</p>
        </div>

        <button
          onClick={handleAccept}
          disabled={isProcessing}
          className="w-full bg-green-600 text-white py-4 rounded-lg font-bold text-lg mb-4 hover:bg-green-700 disabled:bg-gray-400"
        >
          {isProcessing ? 'Processando...' : 'SIM, ADICIONAR AO PEDIDO'}
        </button>

        <button
          onClick={handleDecline}
          className="w-full text-gray-600 py-2 hover:text-gray-800"
        >
          Não, obrigado
        </button>
      </div>

      <div className="text-center bg-yellow-50 py-4 rounded-lg">
        <p className="text-sm text-gray-600 mb-1">⏱️ Oferta expira em</p>
        <p className="text-2xl font-bold">
          {Math.floor(timeLeft / 60)}:{(timeLeft % 60).toString().padStart(2, '0')}
        </p>
      </div>
    </div>
  );
}
```

**Task 2: Upsell Service** (`/backend/src/services/upsell.service.ts`)
```typescript
export class UpsellService {
  async acceptUpsell(sessionId: string, upsellProduct: any): Promise<PaymentResult> {
    const session = await sessionService.getSession(sessionId);
    
    if (!session || session.state !== 'purchased') {
      throw new Error('Invalid session for upsell');
    }

    const paymentData = await sessionService.getPaymentData(sessionId);
    
    // Create new PIX payment for upsell
    const result = await paymentService.createPixPayment({
      amount: upsellProduct.price,
      description: upsellProduct.name,
      payer: {
        name: paymentData.customerName,
        email: paymentData.customerEmail,
        phone: paymentData.customerPhone
      },
      metadata: {
        funnelId: session.funnelId,
        sessionId: session.id,
        type: 'upsell'
      }
    });

    // Update session state
    await sessionService.updateSessionState(sessionId, 'upsell_shown');

    return result;
  }

  async declineUpsell(sessionId: string): Promise<void> {
    await sessionService.updateSessionState(sessionId, 'completed');
    
    // Track decline for analytics
    await analyticsService.trackEvent({
      type: 'upsell_declined',
      sessionId,
      timestamp: new Date()
    });
  }

  async processUpsellPayment(paymentId: string): Promise<void> {
    const payment = await mercadopago.payment.get(paymentId);
    const metadata = payment.body.metadata;
    const sessionId = metadata.session_id;

    // Get original order
    const order = await orderService.getBySessionId(sessionId);

    if (!order) {
      throw new Error('Order not found');
    }

    // Add upsell item to order
    await orderService.addUpsellItem(order.id, {
      productId: 'upsell-product-id',
      productName: 'E-book de Receitas',
      amount: payment.body.transaction_amount,
      type: 'upsell'
    });

    // Update session state
    await sessionService.updateSessionState(sessionId, 'completed');
  }
}

export const upsellService = new UpsellService();
```

**Task 3: Upsell Routes** (`/backend/src/routes/upsell.routes.ts`)
```typescript
router.post('/upsell/accept', upsellController.accept);
router.post('/upsell/decline', upsellController.decline);
```

**Task 4: Upsell Controller** (`/backend/src/controllers/upsell.controller.ts`)
```typescript
export class UpsellController {
  async accept(req: Request, res: Response) {
    try {
      const sessionId = req.cookies.funnel_session_id;
      
      // Get upsell product configuration
      const upsellProduct = {
        name: 'E-book de Receitas Saudáveis',
        price: 47.00
      };

      const result = await upsellService.acceptUpsell(sessionId, upsellProduct);
      
      res.json(result);
    } catch (error) {
      console.error('Upsell acceptance error:', error);
      res.status(500).json({ error: 'Failed to process upsell' });
    }
  }

  async decline(req: Request, res: Response) {
    try {
      const sessionId = req.cookies.funnel_session_id;
      await upsellService.declineUpsell(sessionId);
      
      res.json({ success: true });
    } catch (error) {
      console.error('Upsell decline error:', error);
      res.status(500).json({ error: 'Failed to decline upsell' });
    }
  }
}

export const upsellController = new UpsellController();
```

---

## Database Schema Updates

### Orders Table

```sql
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  funnel_id UUID NOT NULL REFERENCES funnels(id),
  session_id VARCHAR(255) NOT NULL,
  payment_id VARCHAR(255) NOT NULL,
  amount DECIMAL(10, 2) NOT NULL,
  status VARCHAR(50) NOT NULL, -- pending, paid, refunded
  customer_name VARCHAR(255) NOT NULL,
  customer_email VARCHAR(255) NOT NULL,
  customer_phone VARCHAR(50),
  items JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_orders_session_id ON orders(session_id);
CREATE INDEX idx_orders_payment_id ON orders(payment_id);
CREATE INDEX idx_orders_funnel_id ON orders(funnel_id);
```

---

## Testing Strategy

### Unit Tests

**Payment Service:**
```typescript
describe('PaymentService', () => {
  it('should create PIX payment', async () => {
    const result = await paymentService.createPixPayment({
      amount: 197.00,
      description: 'Test Product',
      payer: {
        name: 'John Doe',
        email: 'john@example.com',
        phone: '(11) 99999-9999'
      },
      metadata: {
        funnelId: 'funnel-123',
        sessionId: 'session-123'
      }
    });

    expect(result.paymentId).toBeDefined();
    expect(result.qrCode).toBeDefined();
    expect(result.qrCodeBase64).toBeDefined();
  });

  it('should get payment status', async () => {
    const status = await paymentService.getPaymentStatus('payment-123');
    expect(['pending', 'approved', 'rejected']).toContain(status);
  });
});
```

**Upsell Service:**
```typescript
describe('UpsellService', () => {
  it('should accept upsell and create payment', async () => {
    const result = await upsellService.acceptUpsell('session-123', {
      name: 'Upsell Product',
      price: 47.00
    });

    expect(result.paymentId).toBeDefined();
  });

  it('should decline upsell and update session', async () => {
    await upsellService.declineUpsell('session-123');
    
    const session = await sessionService.getSession('session-123');
    expect(session.state).toBe('completed');
  });
});
```

### Integration Tests

**Complete Payment Flow:**
```typescript
describe('Payment Flow', () => {
  it('should complete full payment flow', async () => {
    // Create payment
    const payment = await request(app)
      .post('/api/payments/create')
      .send({
        amount: 197.00,
        description: 'Test Product',
        payer: {
          name: 'John Doe',
          email: 'john@example.com',
          phone: '(11) 99999-9999'
        }
      })
      .expect(200);

    expect(payment.body.paymentId).toBeDefined();

    // Simulate webhook
    await request(app)
      .post('/webhooks/mercadopago')
      .send({
        type: 'payment',
        data: { id: payment.body.paymentId }
      })
      .expect(200);

    // Check order created
    const order = await orderService.getByPaymentId(payment.body.paymentId);
    expect(order).toBeDefined();
    expect(order.status).toBe('paid');
  });
});
```

### E2E Tests

**Checkout to Upsell Flow:**
```typescript
describe('Checkout E2E', () => {
  it('should complete checkout and upsell', async () => {
    await page.goto('/sales');

    // Fill checkout form
    await page.fill('[name="name"]', 'John Doe');
    await page.fill('[name="email"]', 'john@example.com');
    await page.fill('[name="phone"]', '(11) 99999-9999');

    // Submit
    await page.click('button[type="submit"]');

    // Wait for QR code
    await page.waitForSelector('img[alt="QR Code PIX"]');

    // Simulate payment approval (via webhook)
    await simulatePaymentApproval();

    // Should redirect to upsell
    await page.waitForURL('/upsell');

    // Accept upsell
    await page.click('button:has-text("SIM, ADICIONAR")');

    // Wait for new QR code
    await page.waitForSelector('img[alt="QR Code PIX"]');

    // Simulate upsell payment
    await simulatePaymentApproval();

    // Should redirect to confirmation
    await page.waitForURL('/confirmation');
  });
});
```

---

## Security Considerations

### Webhook Validation

```typescript
function validateMercadoPagoWebhook(req: Request): boolean {
  const signature = req.headers['x-signature'];
  const requestId = req.headers['x-request-id'];
  
  // Validate signature using Mercado Pago's algorithm
  // Implementation depends on Mercado Pago's documentation
  
  return true; // Placeholder
}
```

### Payment Data Encryption

```typescript
import crypto from 'crypto';

function encryptPaymentData(data: any): string {
  const cipher = crypto.createCipher('aes-256-cbc', process.env.ENCRYPTION_KEY!);
  let encrypted = cipher.update(JSON.stringify(data), 'utf8', 'hex');
  encrypted += cipher.final('hex');
  return encrypted;
}

function decryptPaymentData(encrypted: string): any {
  const decipher = crypto.createDecipher('aes-256-cbc', process.env.ENCRYPTION_KEY!);
  let decrypted = decipher.update(encrypted, 'hex', 'utf8');
  decrypted += decipher.final('utf8');
  return JSON.parse(decrypted);
}
```

---

## Deliverables

### Week 5
- [ ] Mercado Pago integration
- [ ] Payment service with PIX support
- [ ] Checkout inline component
- [ ] PIX payment display with QR code
- [ ] Payment polling mechanism
- [ ] Webhook handler
- [ ] Order management system

### Week 6
- [ ] Upsell page component
- [ ] Upsell service
- [ ] One-click upsell flow
- [ ] Countdown timer
- [ ] Payment data encryption
- [ ] E2E tests for complete flow

---

## Acceptance Criteria

- [ ] User can complete checkout with PIX
- [ ] QR code displays correctly
- [ ] Payment confirmation happens within 5 seconds of PIX payment
- [ ] User redirected to upsell after main purchase
- [ ] Upsell acceptance generates new PIX QR code
- [ ] Upsell decline redirects to confirmation
- [ ] Timer expires after 10 minutes
- [ ] All payment data encrypted
- [ ] Webhooks processed correctly
- [ ] Orders tracked in database

---

## Next Phase

Proceed to **[Phase 4: Analytics & Polish](./phase-4-analytics.md)** after completing all deliverables and tests.

---

**Phase Owner:** Full-Stack Team  
**Last Updated:** January 3, 2026
