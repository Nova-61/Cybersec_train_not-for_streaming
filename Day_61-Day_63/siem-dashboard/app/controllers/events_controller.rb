class EventsController < ApplicationController
  # GET /events
  def index
    @events = Event.all.order(event_time: :desc)
  end

  # GET /events/:id
  def show
    @event = Event.find(params[:id])
  end

  # GET /events/new
  def new
    @event = Event.new
  end

  # POST /events
  def create
    @event = Event.new(event_params)

    if @event.save
      redirect_to @event, notice: "Событие создано!"
    else
      render :new, status: :unprocessable_entity
    end
  end

  # GET /events/:id/edit
  def edit
    @event = Event.find(params[:id])
  end

  # PATCH/PUT /events/:id
  def update
    @event = Event.find(params[:id])

    if @event.update(event_params)
      redirect_to @event, notice: "Событие обновлено!"
    else
      render :edit, status: :unprocessable_entity
    end
  end

  # DELETE /events/:id
  def destroy
    @event = Event.find(params[:id])
    @event.destroy
    redirect_to events_path, notice: "Событие удалено!", status: :see_other
  end

  private

  # Strong Parameters — белый список полей, разрешённых для mass assignment
  def event_params
    params.require(:event).permit(:title, :description, :event_time)
  end
end
