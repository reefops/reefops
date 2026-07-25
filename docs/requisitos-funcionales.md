# ReefOps — Requisitos funcionales

La asignación a dominios, componentes y evidencias se mantiene en la
[matriz de alineación](alineacion-requisitos-arquitectura.md).

## 1. Visión

ReefOps será una aplicación para registrar, supervisar y coordinar el cuidado de
Veril y de otros acuarios de agua dulce, salada o salobre. Debe servir tanto a
una persona aficionada como a un equipo que mantenga varios sistemas.

La aplicación debe poder responder rápidamente:

1. ¿Cómo se encuentra cada sistema?
2. ¿Qué necesita atención ahora o próximamente?
3. ¿Qué cambió antes de que apareciera un problema?
4. ¿Qué seres vivos, equipos, productos y existencias contiene?
5. ¿Qué debe hacer el responsable y cómo debe hacerlo?

## 2. Principios del producto

- El objeto principal no es solo el acuario, sino el **sistema acuático
  completo**.
- Cada dato debe conservar su fecha, autor, método y procedencia.
- Agua dulce, salada y salobre tendrán parámetros y flujos adaptados.
- La información podrá mantenerse privada, compartirse de forma controlada o
  publicarse, sin duplicar los datos.
- Registrar una actividad habitual debe requerir pocos pasos.
- La aplicación ayudará a decidir, pero no presentará una inferencia de IA como
  un diagnóstico veterinario confirmado.
- Se mantendrá una trazabilidad completa de mediciones, dosis, cambios,
  incidencias y responsables.

## 3. Conceptos principales

### 3.1 Instalación

Ubicación física que contiene uno o más sistemas: vivienda, negocio, sala,
laboratorio o instalación pública.

### 3.2 Sistema acuático

Conjunto de recipientes y equipos conectados hidráulicamente. Puede incluir:

- urna principal;
- sump;
- refugio;
- depósito de reposición;
- depósito de preparación de agua;
- depósito de cambio de agua;
- cuarentena u hospital;
- tanque de esquejes, reproducción o cría;
- rebosadero y cajas externas;
- reactores y filtros;
- circuitos auxiliares.

Dos urnas visualmente separadas podrán pertenecer al mismo sistema si comparten
agua. El volumen del sistema se calculará sin confundir volumen nominal,
volumen útil y volumen neto estimado.

### 3.3 Componente

Recipiente, equipo, sensor, circuito o elemento que forma parte de un sistema.

### 3.4 Organismo

Ser vivo individual o colonia/grupo gestionado dentro del sistema: pez, coral,
planta, invertebrado, alga, microorganismo cultivado u otro.

## 4. Gestión de instalaciones y sistemas

### RF-001. Instalaciones

- Crear y gestionar varias instalaciones.
- Guardar dirección, zona horaria, contactos y notas de acceso.
- Agrupar los sistemas por sala, cliente, proyecto o finalidad.
- Mostrar un resumen de salud y alertas por instalación.

### RF-002. Ficha del sistema acuático

- Nombre, fotografía, descripción y estado.
- Tipo: dulce, marino, arrecife, salobre, plantado, gambario, biotopo,
  cuarentena, hospital, cría u otro.
- Fecha de montaje, madurez estimada y objetivo del sistema.
- Dimensiones, volumen nominal, útil y neto.
- Sustrato, roca, decoración y hábitat.
- Unidades, moneda, idioma y zona horaria.
- Rangos objetivo configurables.
- Etiquetas y campos personalizados.
- Archivado del sistema sin eliminar su historial.

### RF-003. Topología hidráulica

- Representar visualmente qué componentes están conectados.
- Registrar dirección del flujo entre urna, sump, refugio, reactores y
  depósitos.
- Indicar volúmenes y caudales por tramo.
- Registrar válvulas, derivaciones, retornos y rebosaderos.
- Identificar qué componentes comparten agua.
- Calcular el volumen total y el volumen afectado por una dosis.
- Mantener versiones de la topología cuando se modifique la instalación.

### RF-004. Gestión del sump

- Configurar cámaras y nivel de operación de cada una.
- Registrar skimmer, calentadores, sondas, bombas, reactores y filtración
  ubicados en cada cámara.
- Definir niveles mínimo, normal, máximo y nivel con la bomba de retorno
  apagada.
- Registrar volumen libre de seguridad frente a un corte eléctrico.
- Controlar acumulación de detritos y calendario de limpieza.
- Registrar calcetines, roller mats, esponjas y otros medios mecánicos.
- Alertar por nivel anómalo, desbordamiento, pérdida de sifón o fallo de
  retorno cuando existan sensores.

### RF-005. Depósitos y preparación de agua

- Gestionar agua de ósmosis, reposición, salmuera y agua preparada.
- Registrar capacidad, nivel actual, lote y fecha de preparación.
- Guardar sal utilizada, salinidad objetivo, temperatura y parámetros del lote.
- Calcular agua y sal necesarias.
- Registrar transferencias entre depósitos y sistema.
- Avisar de existencias o niveles bajos.

### RF-005A. Ciclo de vida y madurez del sistema

Cada sistema tendrá una etapa operativa explícita:

- diseño;
- montaje;
- llenado y puesta en marcha;
- ciclado;
- maduración;
- introducción progresiva;
- operación estable;
- reconfiguración;
- parada temporal;
- desmontaje y archivado.

ReefOps deberá:

- ofrecer procedimientos y criterios de paso entre etapas;
- registrar el método de ciclado, fuente de bacterias y material maduro;
- seguir amoníaco, nitrito, nitrato, alcalinidad y otros indicadores relevantes;
- estimar el estado del biofiltro sin presentar la estimación como una
  medición directa;
- impedir que una fecha fija sustituya a la comprobación mediante datos;
- planificar incorporaciones graduales y su impacto esperado;
- identificar periodos de estabilización tras cambios importantes;
- detectar patrones compatibles con síndrome del acuario nuevo o antiguo;
- conservar hitos, cambios de etapa y evidencias;
- generar una lista de comprobación antes de declarar el sistema preparado;
- registrar el desmontaje, destino de habitantes y eliminación de materiales.

### RF-006. Modelo espacial de urnas y recipientes

Cada urna, sump, refugio, depósito y recipiente relevante podrá tener un modelo
espacial propio. El modelo combinará geometría, zonas, superficies y elementos
posicionados.

- Definir forma rectangular, cilíndrica, panorámica, esquinera o personalizada.
- Registrar dimensiones interiores, nivel de agua y grosor del sustrato.
- Representar cristales, fondo, superficie, rebosadero, divisores y cámaras.
- Dibujar roca, madera, sustrato, estructuras, macetas y soportes.
- Diferenciar volumen de agua libre, volumen ocupado y zonas inaccesibles.
- Definir zonas con nombre: frente, fondo, isla izquierda, cámara de retorno,
  refugio u otras.
- Mantener una escala real en centímetros o milímetros.
- Guardar versiones del diseño para consultar cómo estaba montado en una fecha.
- Importar opcionalmente un modelo o plano y exportar una representación
  estructurada.

El modelo espacial será independiente de la fotografía de portada y no exigirá
un escaneo 3D para funcionar.

### RF-007. Posicionamiento 3D mediante vistas 2D

La edición principal se realizará con vistas 2D sincronizadas:

- planta, vista desde arriba, para los ejes horizontal y profundidad;
- alzado frontal para horizontal y altura;
- alzado lateral para profundidad y altura;
- secciones o vistas personalizadas en recipientes complejos.

Al colocar o mover un elemento en dos vistas compatibles, ReefOps calculará una
posición tridimensional `(x, y, z)`. La interfaz deberá:

- mostrar guías, rejilla, escala, cotas y ajuste a superficies;
- permitir introducir coordenadas y dimensiones manualmente;
- arrastrar elementos entre zonas;
- seleccionar punto de anclaje y orientación;
- indicar si el elemento está sobre el sustrato, fijado a una roca, suspendido,
  flotante o unido a un cristal;
- comprobar que la posición se encuentra dentro del recipiente;
- advertir de solapes imposibles y permitir solapes intencionados;
- ofrecer una previsualización 3D navegable, aunque la edición se realice en 2D;
- funcionar de manera razonable en móvil y ofrecer mayor precisión en tableta o
  escritorio.

También podrá obtenerse una posición aproximada marcando el mismo elemento en
dos fotografías tomadas desde ángulos conocidos. Toda posición indicará su
precisión: exacta, estimada o aproximada.

### RF-008. Elementos posicionables

Podrán ubicarse en el espacio:

- peces territoriales o refugios habituales, cuando tenga sentido;
- plantas individuales, macizos y zonas plantadas;
- corales individuales, colonias y esquejes;
- anémonas, moluscos y otros invertebrados;
- roca, madera, cuevas, arena y decoración;
- bombas, retornos, rebosaderos, tomas y salidas;
- calentadores, sondas, difusores, skimmer y filtración;
- luces y módulos de iluminación;
- comederos, dosificadores y puntos de adición;
- sensores y puntos de muestreo;
- zonas de mantenimiento o acceso que deban mantenerse libres.

Cada instancia posicionada podrá tener:

- geometría simple, contorno 2D o volumen aproximado;
- dimensiones y radio de crecimiento esperado;
- orientación e inclinación;
- ficha, fotografías y cronología propias;
- fecha desde la que ocupa esa posición;
- historial de movimientos;
- relaciones de soporte, proximidad, sombra o contacto con otros elementos.

Los organismos que se desplazan libremente podrán asociarse a una zona
preferente, territorio, ruta o intervalo de profundidad en vez de a un punto
fijo.

### RF-008A. Zonas y máscaras espaciales

Un elemento individual o conjunto podrá representarse mediante una **máscara**,
es decir, una región que indique el espacio que ocupa en vez de reducirlo a un
punto o icono. Será aplicable, entre otros casos, a:

- colonias de coral y conjuntos de esquejes;
- macizos, tapizantes y grupos de plantas;
- praderas, algas y organismos incrustantes;
- arena, roca, madera y otras estructuras;
- territorios, zonas de refugio o alimentación;
- áreas con algas, cianobacterias, plagas, daño o tejido perdido;
- zonas funcionales del sump o áreas reservadas para mantenimiento.

Las máscaras podrán crearse:

- dibujando un polígono;
- mediante pincel y borrador;
- combinando o restando formas geométricas;
- trazando el contorno sobre una fotografía calibrada;
- a partir de una propuesta de segmentación realizada por IA y confirmada por
  el usuario;
- duplicando y ajustando una máscara de una fecha anterior.

Cada máscara tendrá:

- propietario: organismo, grupo, estructura, incidencia o zona;
- recipiente y vista a la que pertenece;
- fecha y hora de validez;
- contorno, superficie, perímetro y escala;
- altura, espesor o intervalo de profundidad opcional;
- categoría, color, etiqueta y nivel de precisión;
- procedencia: manual, importada, calculada o sugerida por IA;
- notas y fotografía de referencia;
- estado de revisión.

### RF-008B. Máscaras 2D vinculadas y volumen 3D

- Mantener máscaras independientes en planta, frontal, lateral y secciones.
- Vincular las máscaras que representan el mismo elemento.
- Proyectar una máscara en las demás vistas como guía.
- Estimar un volumen tridimensional a partir de dos o más máscaras compatibles.
- Permitir corregir manualmente el volumen o limitarlo a una superficie de
  soporte.
- Representar elementos planos, ramificados, esféricos, columnares o de forma
  libre.
- Admitir huecos, islas y regiones desconectadas dentro de una misma máscara.
- Dividir una máscara cuando una colonia o macizo se fragmente.
- Unir máscaras cuando varios elementos pasen a gestionarse como un conjunto.
- Conservar la relación entre conjunto e individuos cuando ambos niveles sean
  relevantes.

El volumen 3D derivado mostrará su incertidumbre y no se considerará una
medición exacta si solo existe una vista o no hay una escala fiable.

### RF-008C. Evolución temporal de zonas

Las máscaras serán versionadas en el tiempo mediante instantáneas fechadas. La
aplicación deberá:

- conservar el contorno de cada observación sin sobrescribir los anteriores;
- mostrar una línea temporal y recuperar el estado espacial de cualquier fecha;
- superponer dos fechas con colores diferentes;
- reproducir la evolución mediante una animación;
- interpolar visualmente entre observaciones, marcando los estados
  interpolados;
- calcular cambios de superficie, perímetro, altura y volumen estimado;
- expresar crecimiento absoluto y porcentual por día, semana o mes;
- detectar expansión, contracción, división, fusión, desplazamiento y pérdida
  de tejido;
- diferenciar crecimiento real de un cambio de encuadre o calibración;
- permitir corregir una observación sin eliminar su versión original;
- relacionar cada cambio con parámetros, luz, flujo, dosis, poda, fragmentación
  o incidencias de la misma época.

La ausencia de una máscara en una fecha no implicará que el organismo haya
desaparecido. Los estados desconocido, no observado y ausente serán distintos.

### RF-008D. Conflictos, proximidad y previsión

- Detectar solapamientos actuales entre máscaras.
- Calcular distancia mínima entre contornos.
- Añadir un margen configurable de seguridad o crecimiento.
- Representar alcance de tentáculos, territorio o zona de influencia como una
  máscara adicional.
- Avisar cuando dos zonas incompatibles entren en contacto o se aproximen al
  margen definido.
- Detectar sombreado y ocupación de espacio por capas.
- Proyectar una máscara futura usando el historial de crecimiento.
- Mostrar la proyección como escenario, nunca como hecho.
- Estimar cuándo podría producirse un contacto si la tendencia continuase.
- Proponer poda, fragmentación o recolocación como opciones sujetas a
  confirmación.

### RF-008E. Capas, edición y visualización

- Organizar máscaras en capas: organismos, estructuras, flujo, iluminación,
  incidencias, territorios y mantenimiento.
- Activar, ocultar, bloquear y ajustar la transparencia de cada capa.
- Filtrar por especie, grupo, estado o fecha.
- Ordenar visualmente regiones superpuestas.
- Seleccionar una zona para abrir su ficha y cronología.
- Comparar vista real, máscara, mapa de calor y proyección futura.
- Mostrar una leyenda y una escala comprensibles.
- Disponer de deshacer/rehacer durante la edición.
- Guardar borradores antes de publicar una nueva observación.
- Facilitar el trazado preciso con ratón, lápiz y pantalla táctil.

### RF-009. Luz, flujo y condiciones espaciales

El modelo podrá representar condiciones que varían dentro de cada recipiente.

**Iluminación**

- Posicionar luminarias con altura, orientación, potencia, canales y ópticas.
- Registrar puntos o mapas de mediciones PAR y lux.
- Interpolar un mapa estimado entre mediciones, mostrando su incertidumbre.
- Representar sombra producida por estructuras y organismos.
- Comparar fotoperiodos y configuraciones de iluminación.

**Flujo**

- Posicionar bombas, retornos, bajadas, rebosaderos y obstáculos.
- Definir caudal, dirección, apertura, patrón y programación de cada fuente.
- Dibujar vectores y trayectorias de flujo en vistas 2D.
- Visualizar una aproximación tridimensional mediante capas o secciones.
- Identificar zonas de flujo bajo, alto, turbulento, directo o variable.
- Registrar mediciones o pruebas observadas, como movimiento de partículas.
- Mostrar cómo cambia el flujo según el modo de bombas o el nivel de agua.
- Representar el recorrido entre cámaras del sump, refugios y reactores.

Los mapas calculados serán estimaciones simplificadas, no simulaciones de
dinámica de fluidos certificadas. La aplicación distinguirá entre condiciones
medidas, introducidas por el usuario, interpoladas y simuladas.

### RF-009A. Evaluación de idoneidad de la ubicación

ReefOps podrá valorar si un organismo está bien situado comparando:

- requisitos de su ficha de especie y de la instancia concreta;
- luz, flujo, profundidad y temperatura estimados en su posición;
- sustrato y superficie de anclaje;
- espacio actual y previsto para crecer;
- distancia y compatibilidad con organismos vecinos;
- alcance de tentáculos, agresividad química y territorialidad;
- sombra actual y futura;
- accesibilidad para alimentación y mantenimiento;
- historial de salud, crecimiento y comportamiento;
- parámetros generales del agua y estabilidad reciente.

El resultado incluirá:

- puntuación global y por dimensión;
- clasificación como adecuada, mejorable, inadecuada o sin datos suficientes;
- explicación de los factores favorables y problemáticos;
- procedencia y antigüedad de los datos utilizados;
- nivel de confianza;
- mediciones necesarias para mejorar la evaluación;
- zonas alternativas destacadas en el plano;
- comparación antes de confirmar un traslado.

La valoración no asumirá que el rango genérico de una especie es correcto para
todos sus individuos. El usuario podrá ajustar preferencias observadas y
confirmar si una recolocación mejoró o empeoró su estado.

### RF-009B. Planificación y simulación de cambios

- Probar posiciones sin alterar el diseño operativo.
- Crear escenarios alternativos de aquascaping.
- Simular crecimiento mediante contornos futuros por fecha.
- Detectar contactos, sombreado y falta de espacio futuros.
- Comparar el efecto estimado de mover, añadir o retirar una bomba o luz.
- Buscar zonas candidatas para un nuevo organismo antes de adquirirlo.
- Mostrar conflictos introducidos por un escenario.
- Convertir un escenario aprobado en un plan de tareas.
- Guardar fotografías y mediciones antes y después del cambio.

## 5. Parámetros, mediciones y analítica

### RF-010. Catálogo de parámetros

- Parámetros comunes: temperatura, pH, amoníaco, amonio, nitrito, nitrato,
  fosfato, oxígeno disuelto, turbidez y conductividad.
- Agua dulce: GH, KH, TDS, CO2, hierro, potasio y otros nutrientes.
- Agua marina: salinidad, densidad, alcalinidad, calcio, magnesio, ORP, yodo y
  otros elementos.
- Iluminación: PAR, lux, espectro, fotoperiodo y temperatura de color.
- Parámetros personalizados con unidad, precisión y rango propios.

### RF-011. Registro de mediciones

- Introducción manual rápida.
- Medición individual o panel de varias pruebas.
- Fecha y hora reales de la toma, aunque se introduzca posteriormente.
- Punto de muestreo: urna, sump, depósito u otro componente.
- Método, kit, reactivo, instrumento, lote y fecha de caducidad.
- Fotografía del resultado o de la lectura.
- Notas y nivel de confianza.
- Importación mediante CSV y API.
- Recepción automática desde sensores compatibles.
- Corrección sin perder el valor original ni la auditoría.

### RF-012. Rangos y alertas

- Rangos objetivo distintos por sistema y etapa.
- Umbrales de advertencia y críticos.
- Detección de valores absolutos, cambios bruscos y tendencias.
- Alertas por ausencia de mediciones.
- Confirmación, asignación, silenciamiento y resolución de alertas.
- Modo mantenimiento para evitar falsas alarmas durante una intervención.

### RF-013. Gráficas y correlación

- Histórico por parámetro y comparación entre parámetros.
- Superponer dosis, alimentación, cambios de agua, incorporación de organismos,
  tratamientos e incidencias.
- Comparar periodos y sistemas.
- Mostrar media, mínimo, máximo, variabilidad y velocidad de cambio.
- Anotar eventos directamente en la gráfica.
- Detectar correlaciones como indicios, sin afirmar causalidad.

### RF-014. Calidad y trazabilidad de las mediciones

- Asociar cada medición con método, instrumento, reactivo, lote y operador.
- Registrar fecha de apertura, caducidad y conservación de reactivos.
- Mantener planes e histórico de calibración de sondas e instrumentos.
- Guardar patrones, soluciones de referencia y resultados de calibración.
- Indicar precisión, resolución e incertidumbre conocida.
- Repetir una medición y conservar todas las réplicas.
- Calcular concordancia entre pruebas manuales, sondas y laboratorios.
- Detectar deriva progresiva, saltos imposibles y sensores posiblemente secos o
  contaminados.
- Solicitar confirmación de valores extremos antes de generar recomendaciones
  sensibles.
- Marcar valores como válidos, sospechosos, rechazados o pendientes.
- Corregir un dato conservando el valor original, motivo, autor y fecha.
- Registrar cadena de custodia de muestras externas.
- Importar informes ICP y otros análisis de laboratorio.
- Distinguir resultado inferior al límite de detección de un valor igual a
  cero.
- Mostrar la calidad del dato en gráficas, cálculos y respuestas del agente.

### RF-015. Experimentos y cambios controlados

- Crear un experimento sobre iluminación, flujo, alimento, aditivo,
  mantenimiento o ubicación.
- Definir objetivo, hipótesis, duración y responsable.
- Elegir organismos, zonas y métricas de resultado.
- Registrar un periodo de referencia anterior al cambio.
- Definir qué variable se modifica y qué condiciones deben mantenerse.
- Advertir cuando se cambien varias variables simultáneamente.
- Programar observaciones, fotografías y mediciones.
- Comparar antes, durante y después.
- Relacionar máscaras de crecimiento y comportamiento visual.
- Marcar interferencias, incidencias y datos faltantes.
- Indicar si los datos son insuficientes para concluir.
- Registrar resultado, interpretación y decisión final.
- Convertir un experimento satisfactorio en un procedimiento estable.
- Revertir de forma planificada una configuración cuando el resultado sea
  adverso.

## 6. Seres vivos

### RF-020. Catálogo de especies

- Nombre científico, nombres comunes y clasificación.
- Tipo de organismo y procedencia geográfica.
- Requisitos ambientales y rango recomendado.
- Tamaño adulto, dieta, comportamiento y dificultad.
- Compatibilidad, territorialidad, toxicidad y riesgo para arrecife.
- Necesidades de luz, flujo, sustrato o refugio.
- Fuentes y fecha de revisión de la información.
- Fichas propias cuando una especie no esté catalogada.

### RF-021. Ficha del organismo

- Registro individual, colonia, cardumen, pareja, grupo o población.
- Especie, nombre propio, cantidad y fotografías.
- Sexo, edad estimada, tamaño, peso y etapa vital cuando proceda.
- Origen: criado en cautividad, silvestre, intercambio o desconocido.
- Proveedor, precio, fecha de adquisición y documentos.
- Fecha de entrada, aclimatación, cuarentena y sistema actual.
- Ubicación concreta, zona o territorio dentro del modelo espacial.
- Requisitos espaciales propios y preferencias observadas, que podrán
  prevalecer sobre los valores genéricos de la especie.
- Estado: activo, en cuarentena, trasladado, cedido, vendido, desaparecido o
  fallecido.
- Marcas identificativas, variedad, morfo, generación o linaje.
- Notas de comportamiento, alimentación y relaciones con otros organismos.

### RF-022. Evolución y bienestar

- Cronología individual con fotografías y observaciones.
- Registro de crecimiento, comparativas fotográficas y evolución de máscaras
  espaciales.
- Puntuación o indicadores de condición corporal, coloración, extensión de
  pólipos, apetito y actividad.
- Mudas, desoves, reproducción, gestación y descendencia.
- Árbol de parentesco o linaje cuando sea relevante.
- Registro de bajas con fecha, síntomas y causa confirmada o sospechada.
- Estadísticas de supervivencia sin ocultar individuos archivados.

### RF-023. Compatibilidad y capacidad

- Comprobar incompatibilidades conocidas entre organismos.
- Advertir de requisitos ambientales incompatibles.
- Estimar carga biológica como orientación, explicando sus limitaciones.
- Evaluar espacio, territorialidad y tamaño adulto.
- Simular la incorporación de un organismo antes de añadirlo.
- Evaluar proximidad, contacto, sombreado y alcance agresivo usando posiciones
  y dimensiones espaciales.

### RF-024. Bioseguridad y cuarentena

- Representar qué sistemas comparten agua, aire, utensilios, personal o
  superficies de trabajo.
- Asignar redes, cubos, sifones, pinzas y otros utensilios a un sistema o zona.
- Identificar utensilios mediante etiquetas, colores, QR o NFC.
- Registrar limpieza, desinfección, producto, concentración y tiempo de
  contacto.
- Gestionar recepción, aislamiento, observación, pruebas y liberación.
- Definir duración y protocolo según especie, procedencia y nivel de riesgo.
- Reiniciar o revisar el periodo si se introduce un nuevo organismo.
- Mantener sistemas de cuarentena separados de los de exposición.
- Registrar proveedor, transporte, lote, compañeros de envío y antecedentes
  sanitarios.
- Controlar movimientos de organismos, agua, roca, plantas y equipos.
- Calcular contactos potenciales ante una incidencia sanitaria.
- Bloquear o advertir traslados desde sistemas en observación.
- Gestionar periodos de vacío sanitario y desinfección entre usos.
- Registrar necropsias, muestras y resultados diagnósticos.
- Definir protocolos de acceso para personal, visitantes y profesionales.
- Mostrar un estado de riesgo biológico por sistema, explicando sus factores.
- Evitar recomendar medicación profiláctica sin base profesional o protocolo
  aprobado.

### RF-025. Bienestar animal

Cada organismo o población podrá evaluarse mediante indicadores adaptados a la
especie:

- apetito y acceso efectivo al alimento;
- actividad y patrón de descanso;
- respiración aparente;
- condición corporal, color y estado externo;
- uso del espacio y acceso a refugios;
- comportamiento social y conductas propias de la especie;
- agresión ejercida o recibida;
- respuesta a cuidadores y alimentación;
- estabilidad ambiental;
- signos de dolor, lesión, estrés o enfermedad;
- evolución, reproducción y longevidad.

El sistema deberá:

- combinar observación humana, visión artificial y datos ambientales;
- conservar cada indicador por separado;
- mostrar tendencias y cambios respecto de la línea base individual;
- evitar ocultar un indicador grave dentro de una media global;
- explicar cualquier valoración agregada;
- diferenciar bienestar observado, inferido y desconocido;
- sugerir observaciones adicionales cuando falten datos;
- permitir evaluaciones programadas y revisiones profesionales;
- crear un plan de mejora y comprobar su resultado;
- priorizar prevención y condiciones ambientales antes que tratamientos no
  fundamentados.

### RF-026. Procedencia, trazabilidad y sostenibilidad

- Registrar si un organismo es criado en cautividad, propagado, capturado,
  intercambiado o de origen desconocido.
- Guardar criador, proveedor, país, región, fecha, lote y cadena de custodia.
- Asociar certificados, permisos, facturas y documentación sanitaria.
- Registrar generación, linaje, parentales y propagación.
- Identificar especies protegidas, invasoras o sujetas a restricciones.
- Advertir sobre posibles limitaciones de tenencia, transporte, cesión o
  liberación, remitiendo a fuentes vigentes.
- Mantener un pasaporte exportable del organismo.
- Registrar destino al vender, ceder, trasladar o fallecer.
- Evitar presentar como sostenible un origen no verificado.
- Medir consumo de electricidad, agua, sal y otros recursos.
- Registrar residuos, envases, reactivos y eliminación de medicamentos.
- Estimar impacto por sistema y periodo mostrando los supuestos.
- Seguir producción e intercambio de esquejes, plantas y organismos criados.
- Permitir objetivos de reducción de consumo y residuos.

### RF-027. Transporte, mudanzas y aclimatación

- Planificar traslados dentro de una instalación o entre ubicaciones.
- Crear inventario de organismos, agua, equipos, recipientes y materiales.
- Definir orden de desmontaje, embalaje, transporte y remontaje.
- Asignar recipientes y etiquetas a organismos o grupos.
- Estimar duración, temperatura, oxigenación y autonomía.
- Registrar parámetros de origen, transporte y destino.
- Mantener una checklist y responsables por etapa.
- Controlar ubicación y estado de cada contenedor.
- Registrar incidencias y evidencias durante el trayecto.
- Guiar aclimatación y evitar aplicar un único método a todas las especies.
- Programar observaciones y mediciones posteriores.
- Comparar estado anterior y posterior.
- Mantener un modo de emergencia si el transporte se retrasa.
- Archivar el plan como plantilla reutilizable.

## 7. Alimentación, aditivos y tratamientos

### RF-030. Alimentos

- Catálogo de alimentos con tipo, marca, composición y formato.
- Especies o grupos destinatarios.
- Cantidad y frecuencia recomendada.
- Lote, apertura, caducidad y existencias.
- Plan de alimentación por sistema u organismo.
- Registro rápido de lo realmente suministrado.
- Ayuno programado y excepciones.
- Seguimiento de aceptación, rechazo y reacción.

### RF-030A. Nutrición avanzada

- Definir necesidades por especie, tamaño, etapa vital y estado reproductivo.
- Registrar tamaño de partícula, textura y forma de suministro.
- Gestionar alimento seco, fresco, congelado, vivo y cultivado.
- Documentar descongelación, lavado, enriquecimiento y preparación.
- Crear rotaciones de dieta y objetivos nutricionales.
- Registrar qué individuos accedieron realmente al alimento.
- Detectar competencia, rechazo, sobrealimentación y organismos desplazados.
- Estimar alimento no consumido y retirado.
- Relacionar dieta con crecimiento, condición corporal y nutrientes del agua.
- Gestionar cultivos auxiliares de fitoplancton, zooplancton, artemia y otros.
- Registrar cepa, lote, densidad, cosecha, alimentación y contaminación de
  cultivos.
- Controlar cadena de frío, apertura y descongelaciones.
- Advertir de ingredientes duplicados entre alimentos y suplementos.
- Permitir revisión nutricional por un profesional.

### RF-031. Aditivos y productos

- Gestionar sales, acondicionadores, fertilizantes, buffers, bacterias,
  suplementos, medicamentos y medios químicos.
- Marca, producto, principio o composición, concentración y unidad.
- Lote, proveedor, coste, fecha de apertura y caducidad.
- Instrucciones, advertencias y condiciones de almacenamiento.
- Sistemas en los que puede utilizarse.
- Objetivo: alcalinidad, calcio, magnesio, nutrientes, trazas, carbono,
  remineralización u otro.
- Inventario por envase y ubicación.
- Avisos por stock bajo, caducidad o lote retirado.

### RF-032. Dosificación

- Dosis manual, planificada, recurrente o automática.
- Producto, cantidad, concentración, vía y punto de adición.
- Volumen de agua sobre el que actúa.
- Responsable y hora prevista/real.
- Dividir una dosis diaria en varias tomas.
- Registrar cambios de concentración o de producto.
- Evitar duplicados mediante confirmación de dosis ya administradas.
- Límites máximos configurables y advertencias por incremento brusco.
- Registrar bombas dosificadoras, canales, calibración y depósito asociado.
- Estimar autonomía del depósito.
- Relacionar la dosis con mediciones anteriores y posteriores.

### RF-033. Calculadoras de dosificación

- Calcular cantidad necesaria desde valor actual, objetivo y volumen neto.
- Mostrar supuestos, fórmula, unidades y redondeo.
- Proponer correcciones graduales cuando un cambio rápido pueda ser peligroso.
- No ejecutar automáticamente una recomendación de IA.
- Guardar el cálculo y exigir confirmación antes de convertirlo en una tarea o
  dosis.

### RF-034. Tratamientos y salud

- Registrar síntomas, lesiones, comportamiento y posibles causas.
- Plan de tratamiento con dosis, duración y recordatorios.
- Indicar si el tratamiento se realiza en el sistema principal u hospital.
- Registrar contraindicaciones para especies sensibles e invertebrados.
- Seguimiento diario de respuesta y efectos adversos.
- Periodos de retirada y eliminación segura cuando corresponda.
- Adjuntar indicaciones de un veterinario o especialista.
- Distinguir claramente diagnóstico confirmado de hipótesis.

## 8. Equipos, sensores y automatización

### RF-040. Inventario de equipos

- Tipo, marca, modelo, serie, fotografía y documentación.
- Ubicación física e hidráulica.
- Fecha de compra, instalación, garantía, coste y proveedor.
- Consumo, capacidad y configuración.
- Estado, horas de funcionamiento y vida útil estimada.
- Historial de averías, mantenimientos y repuestos.

### RF-041. Mantenimiento preventivo

- Planes por equipo o componente.
- Instrucciones y checklist.
- Periodicidad por calendario, horas de uso o condición.
- Registro de limpieza, calibración y sustitución.
- Evidencias antes/después.
- Próxima intervención calculada automáticamente.

### RF-042. Sensores e integraciones

- Asociar sensores a un sistema y punto de medición.
- Estado de conexión, batería, calibración y última lectura.
- Conservar datos sin conexión y sincronizarlos posteriormente.
- Detectar lecturas imposibles, sonda seca o sensor obsoleto.
- API e integraciones con controladores, domótica y servicios externos.
- Importar datos sin quedar ligado a un único fabricante.

### RF-043. Control y automatización

- Programar luces, bombas, calentadores, enfriadores, skimmer, alimentación,
  reposición y dosificación.
- Escenas de alimentación, mantenimiento, noche, tormenta o emergencia.
- Reglas con condiciones, retardos, histéresis y límites de seguridad.
- Registro de cada orden, resultado y anulación manual.
- Modo manual con retorno controlado al modo automático.
- Parada segura ante fuga, sobretemperatura, bajo nivel o fallo de bomba.
- Las automatizaciones críticas deberán funcionar localmente cuando sea
  posible, incluso sin Internet.

## 9. Bitácora y operación

### RF-050. Bitácora unificada

- Cronología por sistema, componente y organismo.
- Entradas de texto, voz, fotografía, vídeo y documentos.
- Eventos automáticos y entradas manuales.
- Categorías: observación, medición, mantenimiento, alimentación, dosis,
  incorporación, traslado, reproducción, incidencia y baja.
- Etiquetas, menciones y enlaces entre entradas.
- Búsqueda por fecha, autor, categoría, organismo, producto o parámetro.
- Comparación fotográfica y vista de calendario.
- Exportación completa.

### RF-051. Tareas y procedimientos

- Tareas únicas y recurrentes.
- Plantillas por tipo de acuario y procedimiento.
- Checklist, instrucciones, materiales y tiempo estimado.
- Responsable, suplente, prioridad y fecha límite.
- Posponer, omitir justificadamente o completar parcialmente.
- Crear tareas automáticamente desde una alerta o un plan.
- Adjuntar evidencia y mediciones requeridas.

### RF-052. Cambios de agua

- Planificar porcentaje o volumen.
- Calcular volumen retirado y añadido.
- Registrar sifonado y componentes limpiados.
- Identificar lote de agua nueva y sus parámetros.
- Registrar temperatura y salinidad antes y después.
- Actualizar existencias y generar una entrada de bitácora.

### RF-053. Incidencias

- Registrar fugas, cortes eléctricos, mortalidad, fallo de equipos,
  contaminación, sobredosificación y otros eventos.
- Gravedad, alcance, inicio, detección y resolución.
- Línea temporal de acciones y responsables.
- Protocolo de respuesta y contactos de emergencia.
- Análisis posterior de causa y medidas preventivas.

### RF-054. Emergencias y resiliencia

Se podrán preparar planes para:

- corte eléctrico;
- fallo de calentador, enfriador, bomba o retorno;
- sobretemperatura o falta de oxígeno;
- fuga, desbordamiento o pérdida de sifón;
- dosificación accidental;
- contaminación química;
- fallo de reposición;
- mortalidad o deterioro simultáneo;
- pérdida de conectividad o control;
- incendio, inundación o evacuación;
- ausencia prolongada del responsable.

Cada plan podrá incluir:

- condiciones de activación;
- gravedad y tiempo estimado antes de situación crítica;
- acciones inmediatas ordenadas;
- acciones prohibidas o de riesgo;
- responsables, suplentes y contactos;
- ubicación de llaves, cuadros, válvulas y material;
- baterías, aireadores, generadores y autonomía estimada;
- organismos o sistemas prioritarios;
- checklist utilizable sin Internet;
- escalado si nadie confirma atención;
- evidencias y mediciones requeridas;
- criterios para dar la emergencia por controlada.

ReefOps también deberá:

- permitir simulacros sin generar alarmas reales;
- revisar periódicamente contactos, equipos y autonomías;
- detectar planes incompletos o caducados;
- activar un panel de emergencia simplificado;
- mantener una línea temporal de decisiones;
- generar un informe posterior y tareas preventivas;
- evitar que el agente ejecute acciones críticas no preautorizadas.

## 10. Asistente con inteligencia artificial

### RF-060. Asistente contextual

- Responder usando los datos autorizados del sistema: mediciones, bitácora,
  organismos, equipos, productos y tareas.
- Resumir el estado y señalar datos faltantes o contradictorios.
- Explicar tendencias y formular hipótesis ordenadas por plausibilidad.
- Proponer comprobaciones y tareas, nunca ejecutarlas sin confirmación.
- Citar las mediciones y entradas en las que basa cada respuesta.
- Recordar el contexto del sistema sin mezclar datos de otros acuarios.

### RF-061. Análisis de imágenes

El usuario podrá subir una o varias imágenes y especificar el sistema, organismo
y fecha. El asistente podrá:

- sugerir posibles especies, variedades o grupos taxonómicos;
- comparar el aspecto de un organismo con fotografías anteriores;
- señalar cambios visibles de color, forma, tamaño o comportamiento aparente;
- localizar posibles algas, plagas, parásitos, lesiones o tejido deteriorado;
- estimar cobertura, crecimiento o extensión de una colonia;
- interpretar, con confirmación humana, pruebas colorimétricas;
- leer etiquetas, lotes, caducidades, pantallas e instrumentos;
- detectar nivel de agua, suciedad visible, burbujas o posibles fugas;
- reconocer elementos ya registrados y proponer su posición en el modelo;
- comparar fotografías sucesivas para estimar desplazamiento o crecimiento;
- proponer máscaras de segmentación para organismos, conjuntos, algas, plagas
  o tejido deteriorado;
- alinear imágenes de fechas diferentes y sugerir cambios de contorno;
- solicitar confirmación del usuario antes de incorporar una máscara o cambio
  a la serie temporal;
- detectar cambios visibles en la dirección del flujo mediante partículas,
  pólipos o plantas, indicando la baja fiabilidad de estas inferencias;
- recomendar nuevas fotografías o mediciones para reducir la incertidumbre;
- convertir hallazgos confirmados por el usuario en una entrada de bitácora.

Cada resultado deberá:

- mostrar nivel de confianza y alternativas;
- diferenciar observación visual de inferencia;
- conservar la imagen original y la versión anotada;
- indicar las limitaciones de iluminación, enfoque, escala y balance de blancos;
- pedir confirmación antes de modificar fichas o crear incidencias;
- evitar afirmar diagnósticos clínicos;
- recomendar atención profesional ante riesgos graves.

### RF-062. Captura guiada

- Guiar al usuario para obtener imágenes comparables.
- Superponer encuadre y escala de una sesión anterior.
- Solicitar vista general y primeros planos.
- Registrar iluminación, distancia y cámara cuando estén disponibles.
- Permitir una tarjeta de color o referencia de tamaño.
- Detectar automáticamente imágenes demasiado borrosas u oscuras.

### RF-063. Ayuda operativa

- Generar borradores de rutinas según el tipo de sistema.
- Explicar paso a paso un procedimiento existente.
- Preparar resúmenes diarios o semanales.
- Sugerir qué medir a continuación y justificarlo.
- Consultar el inventario antes de sugerir un producto.
- Detectar posible duplicidad de productos o ingredientes.
- Traducir y resumir manuales, conservando el enlace al original.

### RF-063A. Asistente conversacional del acuario

Cada sistema podrá disponer de un agente contextual capaz de conversar sobre
sus datos y su historial. El usuario podrá preguntarle, por ejemplo:

- cómo se encuentra el acuario y qué requiere atención;
- qué tareas están pendientes;
- cuándo se realizó el último cambio de agua o mantenimiento;
- cómo han evolucionado los parámetros y organismos;
- qué pudo preceder a una anomalía;
- cómo cuidar un habitante concreto;
- si un organismo nuevo podría ser compatible;
- en qué zona podría situarse una planta o coral;
- qué mediciones conviene realizar antes de tomar una decisión;
- qué productos, alimentos o repuestos están disponibles;
- cómo preparar una intervención o ausencia prolongada.

El agente podrá generar:

- resúmenes diarios, semanales o bajo demanda;
- planes de mantenimiento adaptados;
- checklist de diagnóstico y recopilación de evidencias;
- borradores de tareas, entradas de bitácora e informes;
- comparaciones entre periodos;
- escenarios de incorporación o traslado de organismos;
- explicaciones adaptadas al nivel de experiencia del usuario.

Cada respuesta relevante deberá diferenciar:

- hechos registrados en ReefOps;
- cálculos o relaciones derivadas;
- conocimiento procedente de fuentes externas;
- hipótesis y recomendaciones;
- información ausente o incierta.

El agente citará las mediciones, fichas, imágenes, entradas y fuentes utilizadas.
No deberá mezclar datos de sistemas distintos salvo petición explícita y
autorizada.

### RF-063B. Conocimiento y especialización del agente

La especialización podrá construirse combinando:

- contexto estructurado del sistema y sus componentes;
- historial, bitácora, imágenes y correcciones del usuario;
- catálogo de especies y compatibilidades;
- manuales de equipos y fichas de productos;
- procedimientos internos aprobados;
- bibliografía y fuentes externas seleccionadas;
- preferencias y objetivos definidos por el propietario.

Se priorizará recuperación contextual sobre los datos autorizados frente a
incluir información privada de un acuario en un entrenamiento global. Si se
ofrece ajuste o entrenamiento específico:

- será opcional y requerirá consentimiento;
- se explicará qué datos se utilizan y con qué finalidad;
- estará aislado por propietario u organización;
- podrá desactivarse y eliminarse;
- se conservará la procedencia y versión del conocimiento;
- no convertirá correcciones particulares en recomendaciones universales.

El propietario podrá aprobar fuentes de confianza, bloquear fuentes y definir
instrucciones propias, como rutinas, límites o criterios de cuidado.

### RF-063C. Incorporación de nuevos habitantes

El agente ofrecerá un flujo específico antes de incorporar un organismo:

1. Identificar especie o conjunto de alternativas.
2. Recopilar tamaño, cantidad, origen y necesidades.
3. Comparar requisitos con parámetros actuales y estabilidad histórica.
4. Revisar compatibilidad, territorialidad, dieta y carga biológica.
5. Evaluar espacio, luz, flujo y posibles ubicaciones.
6. Considerar el tamaño adulto y el crecimiento futuro.
7. Revisar cuarentena, aclimatación y riesgos sanitarios.
8. Detectar información insuficiente o contradictoria.
9. Presentar ventajas, riesgos, confianza y condiciones necesarias.
10. Crear, si el usuario lo confirma, una lista de preparación y seguimiento.

El resultado no será un simple “sí/no”. Podrá indicar compatible, compatible
con condiciones, desaconsejado o sin datos suficientes, explicando los motivos.

### RF-063D. Integración con asistentes inteligentes y voz

ReefOps podrá integrarse con asistentes inteligentes, altavoces, domótica,
aplicaciones de voz y automatizadores externos. Permitirá:

- consultar estado, alertas, parámetros y tareas;
- dictar observaciones y crear borradores de bitácora;
- registrar una medición mediante voz;
- marcar una tarea como realizada;
- iniciar un temporizador de alimentación o mantenimiento;
- solicitar un resumen audible;
- recibir anuncios de alertas según su gravedad;
- identificar el sistema cuando el usuario gestione varios;
- continuar una conversación iniciada en otro dispositivo.

Ejemplos:

- “¿Cómo está Veril?”
- “Registra una temperatura de 25,4 grados.”
- “¿Cuándo limpiamos el skimmer por última vez?”
- “Añade a la bitácora que el coral está retraído.”
- “¿Qué debo comprobar antes de introducir este pez?”

Los canales de voz deberán:

- confirmar sistema, valor y unidad cuando exista ambigüedad;
- leer de vuelta cualquier dato antes de guardarlo si el riesgo lo requiere;
- distinguir voces o exigir autenticación para información privada;
- evitar pronunciar datos sensibles en dispositivos compartidos;
- impedir dosificación, medicación, cambios de configuración y control crítico
  sin una confirmación reforzada;
- conservar transcripción, autor y dispositivo de las operaciones aceptadas.

### RF-063E. Servidor MCP de ReefOps

ReefOps podrá ofrecer un servidor compatible con Model Context Protocol (MCP)
para que agentes y herramientas autorizadas consulten o utilicen las
capacidades de la plataforma.

El servidor podrá exponer:

**Recursos**

- ficha y estado de sistemas autorizados;
- topología hidráulica y modelo espacial;
- habitantes y ubicaciones;
- parámetros e históricos;
- bitácora, tareas e incidencias;
- equipos, productos e inventario;
- informes y evidencias visuales;
- procedimientos y documentación.

**Herramientas**

- consultar estado y tendencias;
- buscar en la bitácora;
- calcular volúmenes o dosis sin ejecutarlas;
- crear borradores de observaciones, tareas o informes;
- registrar mediciones con validación;
- proponer cambios o planes;
- solicitar un análisis visual;
- generar una vista compartida, sujeta a permisos y confirmación.

**Prompts o flujos reutilizables**

- revisión diaria;
- análisis de parámetros;
- preparación de un cambio de agua;
- evaluación de un nuevo habitante;
- recopilación de información para veterinario;
- investigación de una alarma;
- parte de mantenimiento.

### RF-063F. Seguridad y gobierno de MCP

- Autenticación por usuario, aplicación o agente.
- Permisos por instalación, sistema, recurso y herramienta.
- Acceso de solo lectura como configuración inicial.
- Tokens revocables, con alcance y caducidad.
- Confirmación humana para escrituras y acciones sensibles.
- Prohibición de control crítico directo mediante herramientas generales.
- Registro de agente, usuario, argumentos, resultado y cambios producidos.
- Límites de frecuencia, volumen de datos y coste.
- Filtrado de datos privados según la vista autorizada.
- Separación entre datos, instrucciones del usuario y contenido externo no
  confiable.
- Protección frente a instrucciones maliciosas contenidas en notas, documentos
  o páginas consultadas.
- Versionado de herramientas y esquemas.
- Entorno de pruebas con sistemas ficticios o instantáneas.

### RF-063G. Agentes de terceros y profesionales

- Conectar agentes externos usando MCP o API.
- Permitir que un veterinario o tienda utilice su propio agente sobre una vista
  compartida.
- Aplicar al agente los mismos límites que a la persona que lo autorizó.
- Mostrar qué agente accedió y qué información consultó.
- Permitir aprobar una sesión, consulta concreta o acceso temporal.
- Recibir propuestas estructuradas sin incorporarlas automáticamente.
- Revocar un agente sin revocar necesariamente el acceso del profesional.
- Marcar claramente las respuestas producidas por un agente de terceros.

### RF-064. Controles de seguridad de IA

- No dosificar, medicar ni accionar equipos sin confirmación explícita.
- No inventar mediciones ni completar datos ausentes silenciosamente.
- Mostrar fuentes internas y, cuando proceda, fuentes externas.
- Marcar contenido generado por IA.
- Permitir corregir el resultado y usar esa corrección como feedback.
- Proteger imágenes y datos privados según permisos.
- Registrar qué recomendación fue aceptada, descartada o modificada.

### RF-065. Cámaras y fuentes de visión

El sistema podrá recibir imágenes o vídeo desde:

- fotografías y vídeos subidos manualmente;
- cámara del teléfono mediante una captura guiada;
- cámaras fijas asociadas a una urna, sump, refugio o depósito;
- cámaras subacuáticas compatibles;
- secuencias temporales o clips generados por dispositivos externos;
- integraciones y API autorizadas.

Cada fuente tendrá:

- nombre, ubicación, orientación y zona cubierta;
- recipiente y vista espacial asociados;
- resolución, frecuencia y capacidades;
- estado de conexión y última captura correcta;
- calendario de actividad;
- zonas privadas o excluidas mediante máscara;
- parámetros conocidos de posición y calibración;
- política de conservación de imágenes y vídeo.

Se podrán asociar varias cámaras a un mismo recipiente para reducir oclusiones y
mejorar la localización espacial. La pérdida de señal, lente sucia, imagen
desenfocada, reflejos excesivos o iluminación insuficiente generarán una
incidencia técnica distinta de una alarma biológica.

### RF-066. Análisis visual programado

El propietario podrá crear planes de análisis con:

- frecuencia configurable en minutos, horas, días o eventos concretos;
- horario y días activos;
- cámara, recipiente y zonas que deben analizarse;
- duración del clip o número de capturas;
- tipos de detección habilitados;
- sensibilidad y umbral mínimo de confianza;
- periodo de confirmación antes de alertar;
- prioridad y destinatarios;
- política de conservación de evidencias.

El análisis podrá ejecutarse:

- periódicamente;
- después de alimentar, dosificar o realizar mantenimiento;
- antes y después de apagar o encender las luces;
- al recibir una alerta de sensores;
- bajo demanda.

Cuando sea posible, las capturas se realizarán bajo condiciones comparables de
ángulo e iluminación. El sistema evitará analizar repetidamente una escena que
no cumpla los requisitos mínimos de calidad.

### RF-067. Detección e inventario visual de habitantes

- Detectar peces, corales, plantas, invertebrados y otros organismos visibles.
- Asociar una detección con una ficha existente.
- Sugerir la creación de una ficha para un habitante no registrado.
- Mantener alternativas cuando la identificación no sea concluyente.
- Contabilizar individuos visibles y expresar la incertidumbre por oclusiones.
- Localizar organismos fijos en el modelo espacial.
- Estimar zonas frecuentes, rutas y territorios de organismos móviles.
- Detectar una posible ausencia sin declararla como baja.
- Reconocer nuevas incorporaciones o posibles organismos no deseados.
- Permitir que el usuario corrija identidad, cantidad y ubicación.
- Conservar el historial de asociaciones y correcciones.

La identificación individual requerirá características visuales suficientes y
no se prometerá cuando varios ejemplares sean indistinguibles.

### RF-068. Línea base visual

Antes de evaluar anomalías, ReefOps construirá una línea base específica para
cada sistema, zona y franja horaria:

- aspecto habitual de organismos y estructuras;
- posiciones y territorios frecuentes;
- patrones normales de actividad, natación, alimentación y descanso;
- coloración y extensión habituales;
- cobertura normal de algas y sustrato;
- claridad del agua, partículas y microburbujas habituales;
- comportamiento del flujo visible;
- variación esperada entre día, noche, alimentación y mantenimiento.

El usuario podrá marcar periodos como normales, anómalos o no utilizables. Una
reconfiguración del aquascaping, iluminación, cámara o población abrirá una
nueva versión de la línea base sin borrar la anterior.

### RF-069. Detección visual de problemas

El sistema podrá buscar indicios de:

**Organismos**

- heridas, manchas, puntos, decoloración, inflamación o deterioro visible;
- pérdida de tejido, blanqueamiento, necrosis o pólipos retraídos;
- hojas dañadas, clorosis, agujeros, derretimiento o crecimiento deficiente;
- respiración aparentemente acelerada;
- natación errática, flotabilidad anómala, inmovilidad o aislamiento;
- rascado, agresión, persecución o territorialidad inusual;
- falta de respuesta durante la alimentación;
- ausencia prolongada respecto del patrón habitual;
- mortalidad potencial, pendiente siempre de confirmación humana.

**Acuario y entorno**

- aparición o expansión de algas, cianobacterias u otras coberturas;
- posibles plagas y organismos no registrados;
- agua turbia, cambio de color, exceso de partículas o microburbujas;
- acumulación de detritos;
- cambios de nivel de agua;
- posible fuga o humedad visible;
- bomba detenida, flujo aparentemente reducido o equipo desplazado;
- coral caído, planta desenterrada o estructura movida;
- cristal, lente o sensor visualmente sucio;
- contacto, sombreado o competencia entre organismos.

Cada detección deberá indicar qué se observó, dónde, cuándo, con qué evidencia y
qué alternativas podrían explicarlo. No se presentará una enfermedad o plaga
como confirmada basándose únicamente en visión artificial.

### RF-069A. Seguimiento temporal y comportamiento

- Seguir objetos visibles entre fotogramas y capturas sucesivas.
- Resumir actividad por zona y periodo.
- Comparar el comportamiento con la línea base del mismo horario.
- Detectar cambios persistentes, no solamente eventos aislados.
- Medir de forma aproximada velocidad, permanencia, rutas y uso del espacio.
- Analizar respuesta a alimentación, iluminación y presencia de otros
  organismos.
- Relacionar anomalías con parámetros, dosis, tareas e incidencias recientes.
- Generar fragmentos breves que muestren el comportamiento que originó una
  detección.
- Evitar métricas de comportamiento cuando la cobertura visual sea
  insuficiente.

### RF-069B. Alarmas generadas por visión

Una detección podrá producir:

1. una observación automática sin notificación;
2. una solicitud de revisión;
3. una advertencia;
4. una alarma crítica.

La decisión tendrá en cuenta confianza, gravedad potencial, persistencia,
extensión, velocidad de cambio y corroboración por otras cámaras o sensores.

Cada alarma incluirá:

- sistema, zona y organismos posiblemente afectados;
- imagen o clip de evidencia;
- contorno, máscara o anotación del hallazgo;
- comparación con una referencia anterior;
- nivel de confianza y gravedad;
- explicación en lenguaje claro;
- datos recientes relacionados;
- comprobaciones recomendadas;
- acciones para confirmar, descartar, silenciar o crear una incidencia.

Se agruparán detecciones repetidas del mismo problema. Una alarma descartada
podrá utilizarse para ajustar el sistema, sin ocultar automáticamente eventos
futuros potencialmente graves.

### RF-069C. Revisión humana y aprendizaje

- Bandeja de detecciones pendientes.
- Confirmar, corregir, descartar o marcar como no concluyente.
- Corregir categoría, organismo, máscara, gravedad y causa.
- Marcar falsos positivos y falsos negativos.
- Añadir una detección omitida sobre la imagen.
- Ajustar sensibilidad por cámara, zona y tipo de problema.
- Crear excepciones temporales durante tratamientos o mantenimiento.
- Medir precisión del sistema a partir de revisiones.
- Registrar versión del modelo y configuración que produjo cada resultado.
- Mantener las correcciones específicas de una instalación separadas de datos
  de otros propietarios.

### RF-069D. Privacidad y operación de la visión

- Definir si el análisis ocurre en el dispositivo, localmente o en la nube.
- Mostrar cuándo una cámara está activa.
- Excluir zonas mediante máscaras de privacidad antes de enviar imágenes.
- Configurar retención distinta para capturas normales y evidencias de alarmas.
- Eliminar audio por defecto salvo necesidad y consentimiento explícitos.
- No exponer cámaras ni capturas mediante una página pública por defecto.
- Permitir compartir únicamente las evidencias seleccionadas con veterinarios,
  tiendas o soporte.
- Cifrar transmisión y almacenamiento.
- Registrar accesos, exportaciones y eliminaciones de material visual.
- Seguir funcionando con análisis limitado o captura diferida durante una
  pérdida de Internet cuando la arquitectura lo permita.

## 11. Colaboración y permisos

### RF-070. Usuarios y roles

- Propietario, administrador, cuidador, técnico, observador e invitado temporal.
- Permisos por instalación, sistema y tipo de operación.
- Acceso temporal para vacaciones o mantenimiento externo.
- Historial de accesos y cambios relevantes.
- Invitados con cuenta y visitantes mediante enlace seguro sin necesidad de
  registrarse.

### RF-071. Coordinación

- Asignación y relevo de tareas.
- Comentarios, menciones y evidencia.
- Registro de quién midió, dosificó o modificó una configuración.
- Confirmación de lectura para avisos críticos.
- Parte de turno y resumen de pendientes.

### RF-071A. Procedimientos, formación y calidad operativa

- Crear procedimientos normalizados versionados.
- Indicar alcance, materiales, riesgos, pasos y criterios de aceptación.
- Registrar autor, revisor, aprobación y fecha de vigencia.
- Exigir lectura o formación antes de asignar ciertas operaciones.
- Mantener competencias y autorizaciones por persona.
- Requerir doble validación para dosis, traslados o cambios de alto riesgo.
- Gestionar relevo de turno y aceptación de pendientes.
- Registrar tiempo empleado, desviaciones y resultado.
- Comparar ejecución real con el procedimiento.
- Crear acciones correctivas y preventivas.
- Realizar auditorías internas con evidencias.
- Gestionar niveles de servicio para clientes.
- Generar informes periódicos de mantenimiento profesional.
- Comparar instalaciones usando métricas normalizadas y permisos adecuados.
- Impedir que una métrica de productividad incentive atajos perjudiciales para
  el bienestar.

## 12. Privacidad, compartición y publicación

### RF-072. Niveles de visibilidad

Cada recurso compatible —sistema, ficha de organismo, medición, gráfica,
bitácora, informe, imagen o documento— tendrá uno de estos niveles:

1. **Privado:** accesible únicamente para el propietario y los miembros
   autorizados.
2. **Compartido:** accesible para personas concretas o mediante una URL segura,
   con alcance y duración definidos.
3. **Público:** visible mediante una página pública y, opcionalmente, indexable
   por buscadores.

La visibilidad se heredará del sistema de forma predeterminada, pero podrá
restringirse en recursos concretos. Un recurso privado nunca se hará visible
por estar enlazado desde otro recurso compartido o público.

### RF-073. Enlaces compartidos

- Generar una URL segura para un sistema completo o una vista seleccionada.
- Compartir sin exigir cuenta cuando el propietario lo permita.
- Proteger el enlace con contraseña o código de acceso.
- Definir fecha y hora de caducidad.
- Limitar el número de accesos cuando sea necesario.
- Revocar el enlace inmediatamente.
- Regenerar la URL si se sospecha que ha sido difundida.
- Permitir o impedir descargas, comentarios y carga de archivos.
- Mostrar al propietario cuándo se accedió por última vez y, cuando sea
  legalmente apropiado, un registro básico de accesos.
- Evitar que una URL compartida aparezca en buscadores.
- Mostrar claramente al visitante qué información está viendo y hasta cuándo
  tiene acceso.

### RF-074. Vistas compartidas

El propietario podrá configurar qué módulos y periodo temporal incluye una
vista, por ejemplo:

- estado actual y alertas;
- parámetros y gráficas;
- organismos seleccionados;
- fotografías y vídeos;
- bitácora o entradas concretas;
- tratamientos y evolución clínica;
- equipos, configuración y diagnósticos técnicos;
- tareas e intervenciones;
- documentos adjuntos.

Se podrán ocultar selectivamente:

- ubicación y datos personales;
- costes, facturas y proveedores;
- credenciales, redes e identificadores de dispositivos;
- notas internas;
- automatizaciones y controles;
- otros sistemas de la misma instalación;
- autores o miembros del equipo.

La vista compartida mostrará datos actualizados o una instantánea inmutable,
según el modo elegido al crearla.

### RF-075. Acceso para profesionales

Se proporcionarán plantillas de acceso:

- **Veterinario:** organismos afectados, historial clínico, alimentación,
  parámetros, tratamientos, imágenes y documentos.
- **Tienda o especialista:** parámetros, equipamiento, productos utilizados,
  bitácora e imágenes seleccionadas.
- **Soporte técnico:** modelo y configuración de equipos, telemetría, alarmas,
  incidencias y registros técnicos.
- **Cuidador temporal:** tareas, procedimientos, alimentación, dosis permitidas,
  contactos y protocolo de emergencias.
- **Consulta de solo lectura:** resumen y elementos elegidos, sin posibilidad
  de modificar datos.

El propietario podrá partir de una plantilla y ajustar cada permiso. Los
profesionales con cuenta podrán acceder a los sistemas que les hayan compartido
desde un panel propio.

### RF-076. Interacción de invitados

- Comentarios vinculados a una medición, imagen, organismo o entrada.
- Respuestas en hilos y menciones.
- Carga opcional de informes, recetas, imágenes o archivos.
- Propuestas de tareas, tratamientos o cambios sin aplicarlos directamente.
- Aprobación explícita del propietario antes de incorporar una propuesta a la
  operación del sistema.
- Firma, autor, fecha y estado de revisión de cada aportación.
- Notificación al propietario de nueva actividad.

### RF-077. Página pública del sistema

- URL pública personalizable.
- Nombre, portada, descripción y características generales.
- Habitantes que el propietario decida mostrar.
- Galería, evolución y entradas públicas de la bitácora.
- Parámetros actuales o históricos seleccionados.
- Equipamiento visible opcionalmente.
- Hitos: montaje, incorporaciones, reproducción, premios o cambios importantes.
- Perfil del responsable o proyecto, si se desea.
- Código QR para abrir la página pública.
- Vista previa antes de publicar.
- Activación y desactivación inmediata.
- Opción para permitir o impedir indexación por buscadores.
- Aviso visible de última actualización.

La página pública no mostrará por defecto ubicación exacta, costes, documentos,
incidencias privadas, datos de contacto, controles de equipos ni credenciales.

### RF-078. Publicación desde la bitácora

- Mantener entradas privadas y públicas dentro de una misma bitácora.
- Publicar una entrada existente sin copiarla.
- Crear una versión pública redactada cuando haya información sensible.
- Programar publicación y retirada.
- Previsualizar cómo la verá una persona anónima.
- Compartir la URL de una entrada o galería concreta.
- Conservar el historial de cambios de visibilidad.

### RF-079. Consentimiento, auditoría y seguridad

- Confirmación clara antes de hacer público un recurso.
- Advertencia si una selección contiene datos personales o sensibles.
- Registro de creación, modificación y revocación de accesos.
- Sesiones revocables para usuarios invitados.
- Autenticación reforzada opcional para accesos profesionales.
- Aplicación del mínimo privilegio.
- Eliminar metadatos sensibles de imágenes públicas, incluida la geolocalización.
- No usar recursos privados o compartidos para entrenar sistemas de IA sin un
  consentimiento específico.
- Permitir exportar y eliminar las aportaciones de un invitado conforme a las
  reglas aplicables.
- Página clara de acceso denegado, caducado o revocado, sin revelar que el
  recurso existe.

## 13. Inventario, costes y compras

### RF-080. Inventario

- Productos, repuestos, reactivos, alimentos y consumibles.
- Cantidad, unidad, ubicación, lote, apertura y caducidad.
- Consumo automático desde actividades confirmadas.
- Umbrales mínimos y lista de reposición.
- Historial de movimientos y ajustes.

### RF-081. Costes

- Compras, mantenimiento, electricidad, agua y organismos.
- Categorías y proveedores.
- Coste por instalación, sistema y periodo.
- Estimación de consumo eléctrico por equipo.
- Exportación para análisis externo.

## 14. Informes, importación y portabilidad

### RF-090. Paneles

- Vista global de instalaciones y sistemas.
- Salud, alertas, tareas, mediciones recientes y autonomía de consumibles.
- Panel configurable según rol.
- Indicadores que expliquen por qué muestran cada estado.

### RF-091. Informes

- Resumen diario, semanal y mensual.
- Informe de parámetros, mantenimiento, salud, bajas y costes.
- Informe de una incidencia.
- Informe compartible con cuidador, cliente, tienda o veterinario.
- Creación de una URL segura desde cualquier informe.
- Informe como vista viva o instantánea cerrada en una fecha determinada.
- Exportación a PDF, CSV y formatos estructurados.

### RF-092. Datos y sincronización

- Copia de seguridad y restauración.
- Sincronización entre dispositivos.
- Modo sin conexión para mediciones, tareas y bitácora.
- Resolución visible de conflictos.
- Importación inicial desde hojas de cálculo.
- Exportación completa de los datos del usuario.

## 15. Notificaciones

### RF-100. Notificaciones y escalado

- Notificaciones dentro de la aplicación, push y correo.
- Canales configurables según gravedad.
- Horarios de silencio con excepción para emergencias.
- Escalado a otro responsable si una alerta crítica no se atiende.
- Agrupación para evitar fatiga de notificaciones.
- Cada alerta debe indicar sistema, causa, gravedad y acción sugerida.

## 16. Requisitos no funcionales iniciales

- Diseño móvil primero, con interfaz web adaptable.
- Funcionamiento local como configuración predeterminada.
- El núcleo y las funciones esenciales no dependerán de una conexión a
  Internet.
- Despliegue inicial mediante contenedores.
- Portabilidad entre equipo local, servidor privado y nube.
- Publicación externa mediante proyecciones explícitas, sin exponer la base de
  datos privada ni la red local.
- Registro rápido utilizable con una sola mano.
- Accesibilidad y contraste adecuados.
- Fechas y unidades consistentes y convertibles.
- Auditoría de operaciones sensibles.
- Trazabilidad completa desde solicitud o comando hasta transacción, eventos,
  consumidores, proyecciones, alertas y acciones resultantes.
- Propagación de identificadores de correlación y causación entre dominios,
  trabajadores, dispositivos, agentes y publicación.
- Historial append-only para cambios sensibles; las correcciones crearán nuevas
  versiones sin sobrescribir silenciosamente los valores anteriores.
- Procedencia y linaje de datos derivados, imágenes, máscaras, cálculos,
  recomendaciones e inferencias.
- Ejecución exclusivamente local de modelos de lenguaje, visión, embeddings y
  agentes; ningún dato se enviará a proveedores externos para inferencia.
- Descarga verificable de modelos y actualizaciones sin transferencia de datos
  del acuario.
- Navegación desde cualquier alerta o cambio hacia sus antecedentes y efectos.
- Cifrado en tránsito y protección de datos almacenados.
- Separación estricta entre organizaciones o propietarios.
- Comprobaciones automatizadas para impedir fugas entre vistas privadas,
  compartidas y públicas.
- URLs compartidas con identificadores no predecibles.
- Metadatos sociales y de buscadores configurables para páginas públicas.
- Rendimiento suficiente para históricos largos de sensores.
- API versionada y arquitectura preparada para integraciones.
- Disponibilidad de funciones básicas sin conexión.

## 17. Propuesta de alcance

### MVP

- Usuarios y varios sistemas.
- Ficha del sistema, sump y componentes básicos.
- Etapa del sistema y seguimiento básico del ciclado.
- Parámetros manuales, rangos, gráficas y alertas.
- Trazabilidad de método, instrumento, reactivo y calidad de la medición.
- Fichas de organismos.
- Procedencia, movimientos y cuarentena básica.
- Observaciones de bienestar configurables.
- Productos, aditivos y registro de dosis.
- Alimentación, tareas, cambios de agua y bitácora con fotografías.
- Equipos e inventario básico.
- Panel de estado.
- Plan de emergencia accesible sin conexión.
- Enlaces compartidos de solo lectura, revocables y con caducidad.
- Selección de módulos visibles y ocultación de datos sensibles.
- Exportación y copia de seguridad.

### Segunda etapa

- Colaboración avanzada.
- Acceso profesional, comentarios y carga de documentos.
- Página pública y publicaciones seleccionadas de la bitácora.
- Depósitos, lotes de agua y topología hidráulica visual.
- Modelo espacial con vistas 2D sincronizadas y previsualización 3D.
- Posicionamiento de organismos, equipos, estructuras y puntos de medición.
- Máscaras manuales, capas y series temporales para representar zonas ocupadas.
- Comparación visual y métricas básicas de expansión o contracción.
- Mapas de luz y flujo basados en datos introducidos y mediciones.
- Tratamientos, reproducción y linajes.
- Bioseguridad avanzada, contactos y trazabilidad sanitaria.
- Nutrición avanzada y cultivos auxiliares.
- Transporte, mudanzas y aclimatación.
- Bienestar longitudinal y planes de mejora.
- Experimentos y evaluación antes/después.
- Procedimientos versionados y formación.
- Analítica y correlaciones.
- Asistente contextual y análisis de imágenes con confirmación humana.
- Agente conversacional con respuestas fundamentadas y flujos para nuevos
  habitantes.
- Integración inicial de voz y asistentes inteligentes.
- Importadores e integraciones iniciales.

### Etapa avanzada

- Sensores en tiempo real.
- Controladores y automatizaciones.
- Evaluación espacial de idoneidad y sugerencias de recolocación.
- Simulación de crecimiento, sombra y escenarios de flujo.
- Reconstrucción espacial asistida mediante fotografías.
- Segmentación asistida por IA, estimación volumétrica y proyección de máscaras.
- Cámaras fijas, análisis visual programado e inventario automático.
- Línea base por sistema y detección persistente de anomalías.
- Alarmas visuales revisables y seguimiento básico de comportamiento.
- Servidor MCP, agentes profesionales y automatizaciones conversacionales.
- Trazabilidad ampliada, sostenibilidad y pasaporte de organismos.
- Comparación profesional entre instalaciones y auditorías.
- Detección predictiva de anomalías.
- Visión temporal para crecimiento y estado de organismos.
- Operación profesional multiinstalación.

## 18. Decisiones pendientes

- Público inicial: aficionado, cuidador profesional o ambos.
- Plataformas de lanzamiento.
- Fabricantes y controladores prioritarios.
- Nivel de precisión esperado para el modelo espacial inicial.
- Formatos de importación y exportación de geometría.
- Dimensionamiento del cálculo local para simulaciones de luz y flujo.
- Arquitectura y coste del procesamiento continuo de vídeo.
- Retención predeterminada de capturas normales y evidencias.
- Problemas visuales prioritarios para entrenar y validar el primer modelo.
- Asistentes domésticos y plataformas de automatización prioritarios.
- Herramientas de escritura que expondrá la primera versión del servidor MCP.
- Fuentes externas aprobadas y política de actualización del conocimiento.
- Necesidad real de ajuste de modelos frente a recuperación contextual.
- Grado de automatización permitido.
- Fuente y licencia del catálogo de especies.
- Modelo de indicadores de bienestar por grupo de organismos.
- Protocolos de cuarentena incluidos y autoridad que podrá aprobarlos.
- Reglas legales y territoriales para especies protegidas o invasoras.
- Nivel de trazabilidad exigido a proveedores y organismos.
- Métodos aceptados para calcular consumo e impacto ambiental.
- Calidad mínima de una medición para alimentar recomendaciones o
  automatizaciones.
- Modelo de suscripción y límites por plan.
- Tratamiento legal y privacidad de imágenes y recomendaciones.
- Datos que nunca podrán hacerse públicos.
- Necesidad de cuenta para comentar o adjuntar documentación.
- Retención y detalle del registro de accesos compartidos.
- Funciones que deberán operar sin Internet durante una emergencia.
- Operaciones que requerirán doble validación.
- Métrica usada para el estado general del sistema.

## 19. Referencias de dominio iniciales

Estas fuentes sirven como punto de partida para diseñar protocolos y criterios.
No sustituyen la revisión profesional ni la adaptación a especie y
jurisdicción:

- [Management of Aquarium Fish — Merck Veterinary Manual](https://www.merckvetmanual.com/exotic-and-laboratory-animals/aquarium-fish/management-of-aquarium-fish)
- [Routine Health Care of Fish — Merck Veterinary Manual](https://www.merckvetmanual.com/all-other-pets/fish/routine-health-care-of-fish)
- [Biosecurity — Ornamental Aquatic Trade Association](https://ornamentalfish.org/what-we-do/set-standards/biosecurity/)
- [Animal Care & Management — Association of Zoos & Aquariums](https://www.aza.org/animal-care-management)

Los contenidos derivados deberán guardar fuente, versión, fecha de revisión y
ámbito de aplicación.
