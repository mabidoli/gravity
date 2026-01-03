# Phase 2: No-Code Editor

**Sprint:** 2 (Weeks 3-4)  
**Focus:** Drag-and-drop page builder with mobile-first design  
**Status:** Ready for Implementation

---

## Overview

Phase 2 delivers the core value proposition of Funil Rápido: a visual, no-code editor that enables users to create professional funnel pages without writing code. The editor is built with a mobile-first approach, ensuring optimal experience on the devices where 75% of Brazilian traffic originates.

---

## Objectives

1. Build drag-and-drop interface with WYSIWYG editing
2. Create comprehensive component library
3. Implement mobile-first design system
4. Develop template system for rapid funnel creation
5. Enable real-time preview and publishing

---

## Features

### Feature 1: Editor No-Code Mobile-First

**Priority:** CRITICAL  
**Complexity:** High  
**Estimated Time:** 2-3 weeks

#### Description

Visual page builder with drag-and-drop functionality, designed mobile-first. Users can create Landing Pages, Sales Pages, Upsell Pages, and Confirmation Pages by dragging components onto a canvas and configuring properties.

#### Key Requirements

**Editor Interface:**
- WYSIWYG canvas with real-time preview
- Component library sidebar
- Properties panel for selected component
- Device preview switcher (mobile/desktop)
- Auto-save every 30 seconds
- Undo/Redo (50 actions)

**Component Library:**
- Text (Heading, Paragraph, List)
- Media (Image, Video)
- Buttons (Primary, Secondary, Link)
- Forms (Text, Email, Checkbox)
- Layout (Section, Column, Spacer)
- Social Proof (Testimonial, Badge)
- Special (Checkout Inline, Upsell Button)

**Mobile-First Design:**
- Mobile viewport (320-767px) as primary design target
- Desktop (768px+) as enhancement
- Touch-optimized (44x44px minimum)
- Responsive grid system

#### Architecture

**Component Structure:**
```
/frontend/src/editor/
├── components/
│   ├── Editor.tsx              # Main editor container
│   ├── Canvas.tsx              # Drag-and-drop canvas
│   ├── ComponentLibrary.tsx    # Sidebar with components
│   ├── PropertiesPanel.tsx     # Component properties editor
│   ├── Toolbar.tsx             # Top bar (save, preview, publish)
│   └── DeviceSwitcher.tsx      # Mobile/desktop toggle
├── components-library/
│   ├── Text/
│   │   ├── Heading.tsx
│   │   ├── Paragraph.tsx
│   │   └── List.tsx
│   ├── Media/
│   │   ├── Image.tsx
│   │   └── Video.tsx
│   ├── Buttons/
│   │   ├── PrimaryButton.tsx
│   │   ├── SecondaryButton.tsx
│   │   └── LinkButton.tsx
│   ├── Forms/
│   │   ├── TextField.tsx
│   │   ├── EmailField.tsx
│   │   └── Checkbox.tsx
│   ├── Layout/
│   │   ├── Section.tsx
│   │   ├── Column.tsx
│   │   └── Spacer.tsx
│   ├── SocialProof/
│   │   ├── Testimonial.tsx
│   │   └── Badge.tsx
│   └── Special/
│       ├── CheckoutInline.tsx
│       └── UpsellButton.tsx
├── hooks/
│   ├── useEditor.ts            # Editor state management
│   ├── useDragDrop.ts          # Drag-and-drop logic
│   └── useAutoSave.ts          # Auto-save functionality
├── types/
│   ├── component.types.ts      # Component type definitions
│   └── page.types.ts           # Page structure types
└── utils/
    ├── componentRegistry.ts    # Component registration
    └── pageSerializer.ts       # Page JSON serialization
```

#### Implementation Tasks

**1. Core Editor Infrastructure**

**Task 1.1: Editor State Management** (`/frontend/src/editor/hooks/useEditor.ts`)
```typescript
interface EditorState {
  components: Component[];
  selectedComponentId: string | null;
  history: Component[][];
  historyIndex: number;
  isDirty: boolean;
  viewport: 'mobile' | 'desktop';
}

interface Component {
  id: string;
  type: string;
  props: Record<string, any>;
  children?: Component[];
}

export function useEditor() {
  const [state, setState] = useState<EditorState>({
    components: [],
    selectedComponentId: null,
    history: [[]],
    historyIndex: 0,
    isDirty: false,
    viewport: 'mobile'
  });

  const addComponent = (component: Component, index?: number) => {
    // Add component to canvas
    // Update history
    // Mark as dirty
  };

  const updateComponent = (id: string, props: Partial<Component['props']>) => {
    // Update component properties
    // Update history
    // Mark as dirty
  };

  const removeComponent = (id: string) => {
    // Remove component
    // Update history
    // Mark as dirty
  };

  const undo = () => {
    // Revert to previous history state
  };

  const redo = () => {
    // Move forward in history
  };

  const selectComponent = (id: string | null) => {
    setState(prev => ({ ...prev, selectedComponentId: id }));
  };

  const setViewport = (viewport: 'mobile' | 'desktop') => {
    setState(prev => ({ ...prev, viewport }));
  };

  return {
    state,
    addComponent,
    updateComponent,
    removeComponent,
    undo,
    redo,
    selectComponent,
    setViewport
  };
}
```

**Task 1.2: Drag-and-Drop Integration** (`/frontend/src/editor/hooks/useDragDrop.ts`)
```typescript
import { useDrag, useDrop } from 'react-dnd';

export function useDragComponent(component: Component) {
  const [{ isDragging }, drag] = useDrag({
    type: 'COMPONENT',
    item: { component },
    collect: (monitor) => ({
      isDragging: monitor.isDragging()
    })
  });

  return { isDragging, drag };
}

export function useDropZone(onDrop: (component: Component, index: number) => void) {
  const [{ isOver }, drop] = useDrop({
    accept: 'COMPONENT',
    drop: (item: { component: Component }, monitor) => {
      const index = calculateDropIndex(monitor.getClientOffset());
      onDrop(item.component, index);
    },
    collect: (monitor) => ({
      isOver: monitor.isOver()
    })
  });

  return { isOver, drop };
}
```

**Task 1.3: Auto-Save Hook** (`/frontend/src/editor/hooks/useAutoSave.ts`)
```typescript
export function useAutoSave(
  data: any,
  saveFn: (data: any) => Promise<void>,
  interval: number = 30000
) {
  const [isSaving, setIsSaving] = useState(false);
  const [lastSaved, setLastSaved] = useState<Date | null>(null);

  useEffect(() => {
    const timer = setInterval(async () => {
      if (data.isDirty) {
        setIsSaving(true);
        try {
          await saveFn(data);
          setLastSaved(new Date());
        } catch (error) {
          console.error('Auto-save failed:', error);
        } finally {
          setIsSaving(false);
        }
      }
    }, interval);

    return () => clearInterval(timer);
  }, [data, saveFn, interval]);

  return { isSaving, lastSaved };
}
```

**2. Component Library Implementation**

**Task 2.1: Base Component Interface** (`/frontend/src/editor/types/component.types.ts`)
```typescript
export interface BaseComponentProps {
  id: string;
  className?: string;
  style?: React.CSSProperties;
  margin?: {
    top?: number;
    right?: number;
    bottom?: number;
    left?: number;
  };
  padding?: {
    top?: number;
    right?: number;
    bottom?: number;
    left?: number;
  };
  background?: string;
  visibility?: {
    mobile?: boolean;
    desktop?: boolean;
  };
}

export interface ComponentDefinition {
  type: string;
  name: string;
  icon: React.ReactNode;
  category: 'text' | 'media' | 'button' | 'form' | 'layout' | 'social-proof' | 'special';
  defaultProps: Record<string, any>;
  propertySchema: PropertySchema[];
  render: (props: any) => React.ReactNode;
}

export interface PropertySchema {
  key: string;
  label: string;
  type: 'text' | 'number' | 'color' | 'select' | 'toggle' | 'image';
  options?: Array<{ label: string; value: any }>;
  defaultValue?: any;
}
```

**Task 2.2: Component Registry** (`/frontend/src/editor/utils/componentRegistry.ts`)
```typescript
import { ComponentDefinition } from '../types/component.types';

class ComponentRegistry {
  private components: Map<string, ComponentDefinition> = new Map();

  register(definition: ComponentDefinition) {
    this.components.set(definition.type, definition);
  }

  get(type: string): ComponentDefinition | undefined {
    return this.components.get(type);
  }

  getByCategory(category: string): ComponentDefinition[] {
    return Array.from(this.components.values())
      .filter(comp => comp.category === category);
  }

  getAll(): ComponentDefinition[] {
    return Array.from(this.components.values());
  }
}

export const componentRegistry = new ComponentRegistry();
```

**Task 2.3: Example Component - Heading** (`/frontend/src/editor/components-library/Text/Heading.tsx`)
```typescript
import { ComponentDefinition } from '../../types/component.types';

export const HeadingComponent: ComponentDefinition = {
  type: 'heading',
  name: 'Título',
  icon: <HeadingIcon />,
  category: 'text',
  defaultProps: {
    text: 'Seu Título Aqui',
    level: 'h1',
    fontSize: 32,
    color: '#000000',
    align: 'left',
    fontWeight: 'bold'
  },
  propertySchema: [
    {
      key: 'text',
      label: 'Texto',
      type: 'text',
      defaultValue: 'Seu Título Aqui'
    },
    {
      key: 'level',
      label: 'Nível',
      type: 'select',
      options: [
        { label: 'H1', value: 'h1' },
        { label: 'H2', value: 'h2' },
        { label: 'H3', value: 'h3' }
      ],
      defaultValue: 'h1'
    },
    {
      key: 'fontSize',
      label: 'Tamanho',
      type: 'number',
      defaultValue: 32
    },
    {
      key: 'color',
      label: 'Cor',
      type: 'color',
      defaultValue: '#000000'
    },
    {
      key: 'align',
      label: 'Alinhamento',
      type: 'select',
      options: [
        { label: 'Esquerda', value: 'left' },
        { label: 'Centro', value: 'center' },
        { label: 'Direita', value: 'right' }
      ],
      defaultValue: 'left'
    }
  ],
  render: (props) => {
    const Tag = props.level as keyof JSX.IntrinsicElements;
    return (
      <Tag
        style={{
          fontSize: `${props.fontSize}px`,
          color: props.color,
          textAlign: props.align,
          fontWeight: props.fontWeight,
          ...props.style
        }}
        className={props.className}
      >
        {props.text}
      </Tag>
    );
  }
};

// Register component
componentRegistry.register(HeadingComponent);
```

**Task 2.4: Implement All Components**
- Text: Heading, Paragraph, List (3 components)
- Media: Image, Video (2 components)
- Buttons: Primary, Secondary, Link (3 components)
- Forms: TextField, EmailField, Checkbox (3 components)
- Layout: Section, Column, Spacer (3 components)
- Social Proof: Testimonial, Badge (2 components)
- Special: CheckoutInline, UpsellButton (2 components)

**Total: 18 components**

**3. Editor UI Components**

**Task 3.1: Main Editor Container** (`/frontend/src/editor/components/Editor.tsx`)
```typescript
export function Editor({ pageId }: { pageId: string }) {
  const editor = useEditor();
  const { isSaving, lastSaved } = useAutoSave(
    editor.state,
    async (state) => {
      await fetch(`/api/pages/${pageId}`, {
        method: 'PUT',
        body: JSON.stringify({ content: state.components }),
        headers: { 'Content-Type': 'application/json' }
      });
    }
  );

  return (
    <div className="flex h-screen">
      <ComponentLibrary onAddComponent={editor.addComponent} />
      
      <div className="flex-1 flex flex-col">
        <Toolbar
          onSave={handleSave}
          onPublish={handlePublish}
          onUndo={editor.undo}
          onRedo={editor.redo}
          isSaving={isSaving}
          lastSaved={lastSaved}
        />
        
        <DeviceSwitcher
          viewport={editor.state.viewport}
          onChange={editor.setViewport}
        />
        
        <Canvas
          components={editor.state.components}
          viewport={editor.state.viewport}
          selectedId={editor.state.selectedComponentId}
          onSelect={editor.selectComponent}
          onUpdate={editor.updateComponent}
          onRemove={editor.removeComponent}
        />
      </div>
      
      {editor.state.selectedComponentId && (
        <PropertiesPanel
          component={getComponentById(editor.state.selectedComponentId)}
          onChange={(props) => editor.updateComponent(editor.state.selectedComponentId!, props)}
        />
      )}
    </div>
  );
}
```

**Task 3.2: Canvas Component** (`/frontend/src/editor/components/Canvas.tsx`)
```typescript
export function Canvas({
  components,
  viewport,
  selectedId,
  onSelect,
  onUpdate,
  onRemove
}: CanvasProps) {
  const { isOver, drop } = useDropZone((component, index) => {
    onUpdate(component, index);
  });

  const canvasWidth = viewport === 'mobile' ? '375px' : '100%';

  return (
    <div
      ref={drop}
      className={`canvas ${isOver ? 'drop-active' : ''}`}
      style={{ width: canvasWidth, margin: '0 auto' }}
    >
      {components.length === 0 ? (
        <div className="empty-state">
          Arraste componentes aqui
        </div>
      ) : (
        components.map((component) => (
          <DraggableComponent
            key={component.id}
            component={component}
            isSelected={component.id === selectedId}
            onSelect={() => onSelect(component.id)}
            onRemove={() => onRemove(component.id)}
          />
        ))
      )}
    </div>
  );
}
```

**Task 3.3: Component Library Sidebar** (`/frontend/src/editor/components/ComponentLibrary.tsx`)
```typescript
export function ComponentLibrary({ onAddComponent }: ComponentLibraryProps) {
  const categories = [
    'text',
    'media',
    'button',
    'form',
    'layout',
    'social-proof',
    'special'
  ];

  return (
    <aside className="w-64 bg-gray-100 border-r overflow-y-auto">
      <div className="p-4">
        <h2 className="text-lg font-bold mb-4">Componentes</h2>
        
        {categories.map((category) => (
          <ComponentCategory
            key={category}
            category={category}
            onAddComponent={onAddComponent}
          />
        ))}
      </div>
    </aside>
  );
}

function ComponentCategory({ category, onAddComponent }: CategoryProps) {
  const components = componentRegistry.getByCategory(category);

  return (
    <div className="mb-6">
      <h3 className="text-sm font-semibold mb-2 text-gray-600">
        {getCategoryLabel(category)}
      </h3>
      
      <div className="space-y-2">
        {components.map((comp) => (
          <ComponentItem
            key={comp.type}
            definition={comp}
            onAdd={() => {
              const newComponent = {
                id: generateId(),
                type: comp.type,
                props: { ...comp.defaultProps }
              };
              onAddComponent(newComponent);
            }}
          />
        ))}
      </div>
    </div>
  );
}
```

**Task 3.4: Properties Panel** (`/frontend/src/editor/components/PropertiesPanel.tsx`)
```typescript
export function PropertiesPanel({ component, onChange }: PropertiesPanelProps) {
  const definition = componentRegistry.get(component.type);
  
  if (!definition) return null;

  return (
    <aside className="w-80 bg-white border-l overflow-y-auto p-4">
      <h2 className="text-lg font-bold mb-4">Propriedades</h2>
      
      {definition.propertySchema.map((property) => (
        <PropertyField
          key={property.key}
          property={property}
          value={component.props[property.key]}
          onChange={(value) => {
            onChange({ [property.key]: value });
          }}
        />
      ))}
      
      <Divider />
      
      <GlobalProperties
        margin={component.props.margin}
        padding={component.props.padding}
        background={component.props.background}
        onChange={onChange}
      />
    </aside>
  );
}
```

**4. Mobile-First Responsive System**

**Task 4.1: Responsive Grid System** (`/frontend/src/editor/styles/grid.css`)
```css
/* Mobile-first grid system */
.grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: 1fr;
}

/* Desktop enhancement */
@media (min-width: 768px) {
  .grid-2 {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .grid-3 {
    grid-template-columns: repeat(3, 1fr);
  }
}

/* Touch optimization */
.touch-target {
  min-width: 44px;
  min-height: 44px;
  padding: 12px;
}

button, a {
  @apply touch-target;
}
```

**Task 4.2: Breakpoint Utilities** (`/frontend/src/editor/utils/breakpoints.ts`)
```typescript
export const breakpoints = {
  mobile: {
    min: 320,
    max: 767
  },
  desktop: {
    min: 768,
    max: Infinity
  }
};

export function shouldShowInViewport(
  component: Component,
  viewport: 'mobile' | 'desktop'
): boolean {
  const visibility = component.props.visibility || { mobile: true, desktop: true };
  return visibility[viewport] ?? true;
}
```

**5. Backend API for Pages**

**Task 5.1: Page Controller** (`/backend/src/controllers/page.controller.ts`)
```typescript
export class PageController {
  async getPage(req: Request, res: Response) {
    const { pageId } = req.params;
    const page = await pageService.getById(pageId);
    res.json(page);
  }

  async updatePage(req: Request, res: Response) {
    const { pageId } = req.params;
    const { content } = req.body;
    
    const updated = await pageService.update(pageId, { content });
    res.json(updated);
  }

  async publishPage(req: Request, res: Response) {
    const { pageId } = req.params;
    await pageService.publish(pageId);
    res.json({ success: true });
  }
}
```

**Task 5.2: Page Service** (`/backend/src/services/page.service.ts`)
```typescript
export class PageService {
  async getById(pageId: string): Promise<FunnelPage> {
    return await db.funnel_pages.findUnique({
      where: { id: pageId }
    });
  }

  async update(pageId: string, data: Partial<FunnelPage>): Promise<FunnelPage> {
    return await db.funnel_pages.update({
      where: { id: pageId },
      data: {
        ...data,
        updated_at: new Date()
      }
    });
  }

  async publish(pageId: string): Promise<void> {
    // Mark page as published
    // Clear cache
    // Generate static HTML if needed
  }
}
```

---

### Feature 7: Sistema de Templates

**Priority:** MEDIUM  
**Complexity:** Medium  
**Estimated Time:** 1 week

#### Description

Pre-built funnel templates to accelerate funnel creation. Users can start with a template and customize it instead of building from scratch.

#### Key Requirements

**Templates (MVP):**
1. **Infoproduto (Curso Online)** - Video sales letter with testimonials
2. **E-commerce (Produto Físico)** - Product gallery with specifications
3. **Serviço Local (Consultoria)** - Professional bio with case studies

**Functionality:**
- Template gallery with previews
- One-click template duplication
- Full customization after selection
- Template metadata (name, description, preview image)

#### Implementation Tasks

**Task 1: Template Data Structure** (`/backend/src/types/template.types.ts`)
```typescript
interface FunnelTemplate {
  id: string;
  name: string;
  description: string;
  category: 'infoproduto' | 'ecommerce' | 'servico';
  previewImage: string;
  pages: {
    landing: PageContent;
    sales: PageContent;
    upsell: PageContent;
    confirmation: PageContent;
  };
}

interface PageContent {
  components: Component[];
  metadata: {
    title: string;
    description: string;
  };
}
```

**Task 2: Template Seeds** (`/backend/src/seeds/templates.seed.ts`)
```typescript
export const templates: FunnelTemplate[] = [
  {
    id: 'infoproduto-curso',
    name: 'Infoproduto - Curso Online',
    description: 'Template para venda de cursos online com vídeo de vendas',
    category: 'infoproduto',
    previewImage: '/templates/infoproduto-preview.jpg',
    pages: {
      landing: {
        components: [
          {
            id: '1',
            type: 'heading',
            props: {
              text: 'Transforme Sua Vida com Nosso Curso',
              level: 'h1',
              fontSize: 36,
              color: '#1a202c',
              align: 'center'
            }
          },
          {
            id: '2',
            type: 'video',
            props: {
              url: 'https://www.youtube.com/embed/dQw4w9WgXcQ',
              width: '100%',
              aspectRatio: '16:9'
            }
          },
          {
            id: '3',
            type: 'button-primary',
            props: {
              text: 'QUERO COMEÇAR AGORA',
              action: 'goto-sales',
              size: 'large',
              backgroundColor: '#48bb78',
              textColor: '#ffffff'
            }
          }
          // ... more components
        ],
        metadata: {
          title: 'Curso Online',
          description: 'Landing page para curso online'
        }
      },
      // ... other pages
    }
  },
  // ... other templates
];
```

**Task 3: Template Service** (`/backend/src/services/template.service.ts`)
```typescript
export class TemplateService {
  async getAll(): Promise<FunnelTemplate[]> {
    return templates;
  }

  async getById(templateId: string): Promise<FunnelTemplate | null> {
    return templates.find(t => t.id === templateId) || null;
  }

  async duplicateToFunnel(templateId: string, userId: string): Promise<Funnel> {
    const template = await this.getById(templateId);
    if (!template) throw new Error('Template not found');

    // Create new funnel
    const funnel = await funnelService.create({
      userId,
      name: `${template.name} - Cópia`,
      slug: generateSlug()
    });

    // Create pages from template
    for (const [pageType, pageContent] of Object.entries(template.pages)) {
      await pageService.create({
        funnelId: funnel.id,
        type: pageType,
        content: pageContent
      });
    }

    return funnel;
  }
}
```

**Task 4: Template Gallery UI** (`/frontend/src/components/TemplateGallery.tsx`)
```typescript
export function TemplateGallery({ onSelectTemplate }: TemplateGalleryProps) {
  const { data: templates, isLoading } = useQuery('templates', fetchTemplates);

  if (isLoading) return <LoadingSpinner />;

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 p-6">
      {templates?.map((template) => (
        <TemplateCard
          key={template.id}
          template={template}
          onSelect={() => onSelectTemplate(template.id)}
        />
      ))}
    </div>
  );
}

function TemplateCard({ template, onSelect }: TemplateCardProps) {
  return (
    <div className="border rounded-lg overflow-hidden hover:shadow-lg transition">
      <img
        src={template.previewImage}
        alt={template.name}
        className="w-full h-48 object-cover"
      />
      
      <div className="p-4">
        <h3 className="font-bold text-lg mb-2">{template.name}</h3>
        <p className="text-gray-600 text-sm mb-4">{template.description}</p>
        
        <button
          onClick={onSelect}
          className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
        >
          Usar Template
        </button>
      </div>
    </div>
  );
}
```

**Task 5: Template Routes** (`/backend/src/routes/template.routes.ts`)
```typescript
router.get('/templates', templateController.getAll);
router.get('/templates/:templateId', templateController.getById);
router.post('/templates/:templateId/duplicate', templateController.duplicate);
```

---

## Testing Strategy

### Unit Tests

**Editor State Management:**
```typescript
describe('useEditor', () => {
  it('should add component to canvas', () => {
    const { result } = renderHook(() => useEditor());
    
    act(() => {
      result.current.addComponent({
        id: '1',
        type: 'heading',
        props: { text: 'Test' }
      });
    });
    
    expect(result.current.state.components).toHaveLength(1);
    expect(result.current.state.isDirty).toBe(true);
  });

  it('should support undo/redo', () => {
    const { result } = renderHook(() => useEditor());
    
    act(() => {
      result.current.addComponent({ id: '1', type: 'heading', props: {} });
      result.current.addComponent({ id: '2', type: 'paragraph', props: {} });
    });
    
    expect(result.current.state.components).toHaveLength(2);
    
    act(() => {
      result.current.undo();
    });
    
    expect(result.current.state.components).toHaveLength(1);
    
    act(() => {
      result.current.redo();
    });
    
    expect(result.current.state.components).toHaveLength(2);
  });
});
```

**Component Registry:**
```typescript
describe('ComponentRegistry', () => {
  it('should register and retrieve components', () => {
    const registry = new ComponentRegistry();
    
    registry.register(HeadingComponent);
    
    const retrieved = registry.get('heading');
    expect(retrieved).toBeDefined();
    expect(retrieved?.name).toBe('Título');
  });

  it('should get components by category', () => {
    const registry = new ComponentRegistry();
    
    registry.register(HeadingComponent);
    registry.register(ParagraphComponent);
    
    const textComponents = registry.getByCategory('text');
    expect(textComponents).toHaveLength(2);
  });
});
```

### Integration Tests

**Page Save/Load:**
```typescript
describe('Page API', () => {
  it('should save page content', async () => {
    const pageContent = {
      components: [
        { id: '1', type: 'heading', props: { text: 'Test' } }
      ]
    };
    
    const response = await request(app)
      .put('/api/pages/test-page-id')
      .send({ content: pageContent })
      .expect(200);
    
    expect(response.body.content).toEqual(pageContent);
  });

  it('should load page content', async () => {
    const response = await request(app)
      .get('/api/pages/test-page-id')
      .expect(200);
    
    expect(response.body.content).toBeDefined();
    expect(response.body.content.components).toBeInstanceOf(Array);
  });
});
```

### E2E Tests

**Complete Editor Flow:**
```typescript
describe('Editor E2E', () => {
  it('should create page from scratch', async () => {
    await page.goto('/editor/new');
    
    // Drag heading component
    await page.dragAndDrop(
      '[data-component="heading"]',
      '.canvas'
    );
    
    // Edit text
    await page.dblclick('.canvas [data-type="heading"]');
    await page.keyboard.type('My Heading');
    
    // Change color
    await page.click('.canvas [data-type="heading"]');
    await page.click('[data-property="color"]');
    await page.fill('[data-property="color"] input', '#ff0000');
    
    // Save
    await page.click('[data-action="save"]');
    await page.waitForSelector('.toast-success');
    
    // Verify saved
    const saved = await page.evaluate(() => {
      return fetch('/api/pages/current').then(r => r.json());
    });
    
    expect(saved.content.components).toHaveLength(1);
    expect(saved.content.components[0].props.text).toBe('My Heading');
    expect(saved.content.components[0].props.color).toBe('#ff0000');
  });

  it('should use template', async () => {
    await page.goto('/templates');
    
    // Select template
    await page.click('[data-template="infoproduto-curso"]');
    
    // Wait for editor to load
    await page.waitForSelector('.canvas');
    
    // Verify components loaded
    const components = await page.$$('.canvas [data-component]');
    expect(components.length).toBeGreaterThan(0);
    
    // Customize
    await page.click('.canvas [data-type="heading"]');
    await page.fill('[data-property="text"] input', 'Meu Curso');
    
    // Publish
    await page.click('[data-action="publish"]');
    await page.waitForSelector('.toast-success');
  });
});
```

---

## Performance Optimization

### Virtualization for Large Component Lists

```typescript
import { FixedSizeList } from 'react-window';

export function ComponentLibrary({ onAddComponent }: ComponentLibraryProps) {
  const allComponents = componentRegistry.getAll();

  return (
    <FixedSizeList
      height={600}
      itemCount={allComponents.length}
      itemSize={60}
      width="100%"
    >
      {({ index, style }) => (
        <div style={style}>
          <ComponentItem
            definition={allComponents[index]}
            onAdd={() => onAddComponent(allComponents[index])}
          />
        </div>
      )}
    </FixedSizeList>
  );
}
```

### Lazy Loading for Components

```typescript
const HeadingComponent = lazy(() => import('./components-library/Text/Heading'));
const ParagraphComponent = lazy(() => import('./components-library/Text/Paragraph'));
// ... etc

export function Canvas({ components }: CanvasProps) {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      {components.map((component) => {
        const Component = getComponentByType(component.type);
        return <Component key={component.id} {...component.props} />;
      })}
    </Suspense>
  );
}
```

---

## Deliverables

### Week 3
- [ ] Editor state management (useEditor hook)
- [ ] Drag-and-drop infrastructure
- [ ] Component registry system
- [ ] 10 basic components (Text, Media, Buttons, Forms)
- [ ] Canvas with device switcher
- [ ] Component library sidebar

### Week 4
- [ ] Properties panel
- [ ] Remaining 8 components (Layout, Social Proof, Special)
- [ ] Auto-save functionality
- [ ] Undo/Redo
- [ ] Template system (3 templates)
- [ ] Template gallery UI
- [ ] E2E tests for editor

---

## Acceptance Criteria

- [ ] User can create blank page and add components
- [ ] Drag-and-drop works smoothly (< 100ms response)
- [ ] Component properties update in real-time
- [ ] Mobile/desktop preview works correctly
- [ ] Auto-save triggers every 30 seconds
- [ ] Undo/Redo works for up to 50 actions
- [ ] User can select and use templates
- [ ] 80% of users can publish first page in < 15 minutes

---

## Next Phase

Proceed to **[Phase 3: Payment & Upsells](./phase-3-payments.md)** after completing all deliverables and tests.

---

**Phase Owner:** Frontend Team  
**Last Updated:** January 3, 2026
